package models

type FootballAPIPlayerInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
	Height      string `json:"height"`
	Weight      string `json:"weight"`
	Photo       string `json:"photo"`
	Position    string `json:"position"`
}

type FootballAPIPlayerResponse struct {
	Response []FootballAPIPlayerEntry `json:"response"`
}

type FootballAPIPlayerEntry struct {
	Statistics []FootballAPIPlayerStatistics `json:"statistics"`
}

type FootballAPIPlayerStatistics struct {
	Games struct {
		Appearances *int    `json:"appearences"`
		Minutes     *int    `json:"minutes"`
		Position    *string `json:"position"`
		Rating      *string `json:"rating"`
		Captain     *bool   `json:"captain"`
	} `json:"games"`

	Shots struct {
		Total *int `json:"total"`
		On    *int `json:"on"`
	} `json:"shots"`

	Goals struct {
		Total    *int `json:"total"`
		Conceded *int `json:"conceded"`
		Assists  *int `json:"assists"`
		Saves    *int `json:"saves"`
	} `json:"goals"`

	Passes struct {
		Total *int `json:"total"`
		Key   *int `json:"key"`
	} `json:"passes"`

	Tackles struct {
		Total         *int `json:"total"`
		Blocks        *int `json:"blocks"`
		Interceptions *int `json:"interceptions"`
	} `json:"tackles"`

	Duels struct {
		Total *int `json:"total"`
		Won   *int `json:"won"`
	} `json:"duels"`

	Dribbles struct {
		Attempts *int `json:"attempts"`
		Success  *int `json:"success"`
		Past     *int `json:"past"`
	} `json:"dribbles"`

	Fouls struct {
		Drawn     *int `json:"drawn"`
		Committed *int `json:"committed"`
	} `json:"fouls"`

	Cards struct {
		Yellow    *int `json:"yellow"`
		YellowRed *int `json:"yellowred"`
		Red       *int `json:"red"`
	} `json:"cards"`

	Penalty struct {
		Won      *int `json:"won"`
		Commited *int `json:"commited"`
		Scored   *int `json:"scored"`
		Missed   *int `json:"missed"`
		Saved    *int `json:"saved"`
	} `json:"penalty"`
}
