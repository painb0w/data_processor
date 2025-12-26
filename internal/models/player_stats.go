package models

import "time"

type PlayerStats struct {
	ID        int
	Name      string
	Position  string
	Stats     any
	UpdatedAt time.Time
}

type AttackerStats struct {
	GamesPlayed     int
	Shots           int
	Goals           int
	Assists         int
	Dribbles        int
	SuccessDribbles int
	FoulsDrawn      int
}

type MidfielderStats struct {
	GamesPlayed int
	Passes      int
	KeyPasses   int
	Assists     int
	Tackles     int
	Dribbles    int
}

type DefenderStats struct {
	GamesPlayed   int
	Tackles       int
	Interceptions int
	Blocks        int
	FoulsCommited int
	DuelsWon      int
	Duels         int
	YellowCards   int
	RedCards      int
}

type GoalkeeperStats struct {
	GamesPlayed    int
	Conceded       int
	Saves          int
	Penalties      int
	PenaltiesSaved int
}
