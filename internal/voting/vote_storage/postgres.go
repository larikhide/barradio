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
	PG *sql.DB
}

// NewPostgresVoteStorage returns new rady to use instance of DB connector
func NewPostgresVoteStorage(url string, mgrt string) (*PostgresVoteStorage, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	vtst := &PostgresVoteStorage{PG: db}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		//log.Println("postgres.WithInstance")
		log.Println(err)
		return vtst, err
	}
	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", mgrt), "barradio", driver)
	if err != nil {
		//log.Println("migrate.NewWithDatabaseInstance")
		log.Println(err)
		driver.Close()
		return vtst, err
	}


	m.Up()
	
	
	//Для автогенерации голосования********************************
	//go vtst.GenerateVotes("3s")
	//*************************************************************
	return vtst, nil
}

// Close disconnects from DB
func (s *PostgresVoteStorage) Close() error {

	return s.PG.Close()
}

func (s *PostgresVoteStorage) GetVoteCategories() ([]string, error) {

	var cat string
	var categories []string
	rows, err := s.PG.Query("select name from category  where is_run = true")
	if err != nil {

		log.Println(err)
		rows.Close()
		return categories, err
	}
	defer rows.Close()
	for rows.Next() {

		err := rows.Scan(&cat)
		if err != nil {
			log.Println(err)
			continue
		}
		categories = append(categories, cat)
	}

	return categories, nil
}

func (s *PostgresVoteStorage) SaveVoteForCategory(vote vote_service.Vote) error {

	// TODO store in DB instead of memory

	var code int
	var category_id uuid.UUID

	row := s.PG.QueryRow("select id, code from category where name = $1 and is_run = true", vote.Category)
	err := row.Scan(&category_id, &code)
	if err != nil {
		return err
	}
	id := uuid.New()
	_, err = s.PG.Exec(`insert into vote (id, category_id , category_code, created_at) values ($1, $2, $3, $4)`, id, category_id, code, vote.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}
func (s *PostgresVoteStorage) GetVotesCountForInterval(start, end time.Time) (map[string]int, error) {

	// TODO fetch from DB

	result := make(map[string]int)
	rows, err := s.PG.Query(`select c.name, v.category_code, count(v.category_code) from vote v  
							inner join category c on c.code = v. category_code
	                         where v.created_at between $1 and $2 
							 group by c.name,v.category_code`, start, end)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	var cat string
	var voteResult int
	var code int

	for rows.Next() {

		err := rows.Scan(&cat, &code, &voteResult)
		if err != nil {
			log.Println(err)
			continue
		}
		result[cat] = voteResult
	}

	return result, nil
}

func (s *PostgresVoteStorage) GenerateVotes(ss string) {

	var vote vote_service.Vote
	var cat string
	for {
		n := (1 + rand.Intn(3-1+1))
		row := s.PG.QueryRow(`select c.name from category c where c.code =$1`, n)
		err := row.Scan(&cat)
		if err != nil {
			return
		}

		vote = vote_service.Vote{
			Category:  cat,
			CreatedAt: time.Now(),
		}
		
		_ = s.SaveVoteForCategory(vote)
		d, _ := time.ParseDuration(ss)
		time.Sleep(d)
	}
}
