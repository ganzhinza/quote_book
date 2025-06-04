package memdb

import (
	"errors"
	"math/rand"
	"quote_book/pkg/entities"
	"quote_book/pkg/utils"
	"sync"
	"time"
)

const garbagePart = 0.1

type safeQuote struct {
	*entities.Quote
	deleted bool
	sync.RWMutex
}

type MemDB struct {
	sync.RWMutex
	idGenerator *utils.IDGenerator
	garbagePart float64
	quotes      map[int]*safeQuote
	authorIndex map[string]map[int]*safeQuote //Maybe better make indexes map[string]map[string][]*entities.Quote for bigger project
	aliveIDs    []int
	aliveIDsMu  sync.Mutex
	deadIDs     map[int]bool
	deadIDsMu   sync.Mutex
}

func New() *MemDB {
	db := &MemDB{
		garbagePart: garbagePart,
		quotes:      make(map[int]*safeQuote),
		authorIndex: make(map[string]map[int]*safeQuote),
		aliveIDs:    make([]int, 0),
		deadIDs:     make(map[int]bool),
		idGenerator: utils.NewIDGenerator(0),
	}
	go db.GarbageCollector()

	return db
}

// блокируем только добавление живых индексов - проблема только если одновременно добавляем в слайс
func (db *MemDB) AddQuote(quote entities.Quote) error {
	if quote.Text == "" {
		return errors.New("blank quote")
	}

	quote.ID = db.idGenerator.GetID()

	sQuote := &safeQuote{Quote: &quote}

	db.quotes[quote.ID] = sQuote
	if db.authorIndex[quote.Author] == nil {
		db.authorIndex[quote.Author] = make(map[int]*safeQuote)
	}
	db.authorIndex[quote.Author][sQuote.ID] = sQuote

	db.aliveIDsMu.Lock()
	db.aliveIDs = append(db.aliveIDs, quote.ID)
	db.aliveIDsMu.Unlock()

	return nil
}

// блокировка на чтение (для работы GC)
func (db *MemDB) GetAllQuotes() ([]entities.Quote, error) {
	db.RLock()
	defer db.RUnlock()
	quotes := make([]entities.Quote, 0, len(db.quotes))

	for _, sQuote := range db.quotes {
		if !sQuote.deleted {
			quotes = append(quotes, *sQuote.Quote)
		}
	}
	return quotes, nil
}

// блокировка на чтение (для работы GC)
func (db *MemDB) GetRandomQuote() (entities.Quote, error) {
	db.RLock()
	defer db.RUnlock()

	id, err := db.GetAliveID()
	if err != nil {
		return entities.Quote{}, err
	}
	return *db.quotes[id].Quote, nil
}

func (db *MemDB) GetAliveID() (int, error) {
	db.RLock()
	defer db.RUnlock()
	id := 0
	for {
		if len(db.aliveIDs)-len(db.deadIDs) == 0 {
			return -1, errors.New("no valid ids")
		}
		id = db.aliveIDs[rand.Intn(len(db.aliveIDs))]
		if !db.quotes[id].deleted {
			return id, nil
		}
	}
}

// логическое удаление, чтобы не останавливать всю базу ради одного удаления
func (db *MemDB) DeleteQuote(id int) error {
	db.RLock()
	defer db.RUnlock()

	_, exists := db.quotes[id]
	if !exists {
		return nil
	}

	db.quotes[id].Lock()
	if !db.quotes[id].deleted {
		db.quotes[id].deleted = true

		db.deadIDsMu.Lock()
		db.deadIDs[id] = true
		db.deadIDsMu.Unlock()
	}

	db.quotes[id].Unlock()

	return nil
}

// блокировка на чтение (для работы GC)
func (db *MemDB) GetAuthorQuotes(author string) ([]entities.Quote, error) {
	db.RLock()
	defer db.RUnlock()

	quotes := make([]entities.Quote, 0, len(db.authorIndex[author]))
	for _, sQuote := range db.authorIndex[author] {
		quotes = append(quotes, *sQuote.Quote)
	}

	return quotes, nil
}

// убираем мусор, когда его больше, чем заданный порог
func (db *MemDB) GarbageCollector() {
	for {
		if float64(len(db.deadIDs))/float64(len(db.quotes)) > db.garbagePart {
			db.Lock() //stop the world

			for id := range db.deadIDs {
				delete(db.authorIndex[db.quotes[id].Author], id)
				delete(db.quotes, id)
			}
			for author := range db.authorIndex {
				if len(db.authorIndex[author]) == 0 {
					delete(db.authorIndex, author)
				}
			}
			db.deadIDs = make(map[int]bool)

			db.aliveIDs = make([]int, 0, len(db.quotes))
			for id := range db.quotes {
				db.aliveIDs = append(db.aliveIDs, id)
			}
			db.Unlock()
		}
		time.Sleep(time.Millisecond * 500)
	}
}
