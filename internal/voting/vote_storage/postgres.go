package vote_storage

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/larikhide/barradio/internal/voting/vote_service"
)

var _ vote_service.VoteStorage = &PostgresVoteStorage{}

// PostgresVoteStorage encapsulates all DB work
type PostgresVoteStorage struct {

	// TODO
	// now data stores in memory - just for tests
	// remove after implementing DB calls
	// votes []struct {
	// 	CreatedAt time.Time
	// 	Category  string
	// }

	PG *sql.DB
}

// NewPostgresVoteStorage returns new rady to use instance of DB connector
func NewPostgresVoteStorage(url string, mgrt string) (*PostgresVoteStorage, error) {
	// dummy to just start app server
	// TODO implement due https://github.com/larikhide/barradio/issues/8
	//return &PostgresVoteStorage{}, nil

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", mgrt), "barradio", driver)
	if err != nil {
		log.Println(err)
		log.Fatalf("cannot run migration: %s", err.Error())
	}
	m.Up()
	if err != nil {
		log.Println(err)
		log.Fatalf("cannot run migration: %s", err.Error())
	}

	vtst := &PostgresVoteStorage{PG: db}

	go vtst.GenerateVotes("3s")

	return vtst, nil
}

// Close disconnects from DB
func (s *PostgresVoteStorage) Close() error {
	// dummy to clean shutdown app server
	// TODO implement due https://github.com/larikhide/barradio/issues/8
	return s.PG.Close()
}

func (s *PostgresVoteStorage) GetVoteCategories() ([]string, error) {

	// TODO make actual DB request

	var cat string
	var categories []string
	rows, err := s.PG.Query("select name from category  where is_run = true")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {

		err := rows.Scan(&cat)
		if err != nil {
			fmt.Println(err)
			continue
		}
		categories = append(categories, cat)
	}

	return categories, nil
}

func (s *PostgresVoteStorage) SaveVoteForCategory(vote vote_service.Vote) error {

	// TODO store in DB instead of memory

	// s.votes = append(s.votes, struct {
	// 	CreatedAt time.Time
	// 	Category  string
	// }{
	// 	CreatedAt: vote.CreatedAt,
	// 	Category:  vote.Category,
	// })

	/*
			id uuid NOT NULL,
		    category_id uuid NOT NULL,
		    category_code integer NOT NULL,
		    created_at timestamp NOT NULL,
	*/
	var code int
	var category_id uuid.UUID

	row := s.PG.QueryRow("select id, code from category where name = $1 and is_run = true", vote.Category)
	err := row.Scan(&category_id, &code)
	if err != nil {
		return err
	}
	id := uuid.New()
	_ , err = s.PG.Exec(`
         insert into vote (id, category_id , category_code,created_at) 
		  values ($1, $2, $3, $4)`, id, category_id, code, vote.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}
func (s *PostgresVoteStorage) GetVotesCountForInterval(start, end time.Time) (map[string]int, error) {

	// TODO fetch from DB

	result := make(map[string]int)

	// for _, vote := range s.votes {
	// 	if vote.CreatedAt.After(start) && vote.CreatedAt.Before(end) {
	// 		result[vote.Category] += 1
	// 	}
	// }

	return result, nil
}

func (s *PostgresVoteStorage) GenerateVotes(ss string) {
	for {
		n := (1 + rand.Intn(3-1+1))
		fmt.Println(n)
		d, _ := time.ParseDuration(ss)
		time.Sleep(d)
	}
}
