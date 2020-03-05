package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	ctoai "github.com/cto-ai/sdk-go"
)

//TriviaJSON is the response body we get from OpenTDB
type TriviaJSON struct {
	ResponseCode int `json:"response_code"`
	Results      []struct {
		Category         string   `json:"category"`
		Type             string   `json:"type"`
		Difficulty       string   `json:"difficulty"`
		Question         string   `json:"question"`
		CorrectAnswer    string   `json:"correct_answer"`
		IncorrectAnswers []string `json:"incorrect_answers"`
	} `json:"results"`
}

//TriviaElement is our internal representation of a question
type TriviaElement struct {
	Question        string
	CorrectAnswer   string
	PossibleAnswers string
}

var client = ctoai.NewClient()

//printWrapper fixes ux.Print, allowing it to print arbitrary, multiple arguments ala fmt.Println
func printWrapper(a ...interface{}) {
	client.Ux.Print(fmt.Sprint(a...))
}

//badErrHandler prints and panics on error
func badErrHandler(err error) {
	if err != nil {
		printWrapper(err)
		panic(err)
	}
}

const diff = "hard"
const amount = 2

func getTrivia(token string) (qa TriviaJSON, err error) {
	reqURL := fmt.Sprintf("https://opentdb.com/api.php?encode=base64&amount=%d&difficulty=%s&type=multiple&token=%s", amount, diff, token)
	resp, err := http.Get(reqURL)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var respJSON TriviaJSON
	err = json.Unmarshal(respBody, &respJSON)
	if err != nil {
		return
	}

	if len(respJSON.Results) == 0 {
		err = fmt.Errorf("No results found")
		return
	}

	// if respJSON.ResponseCode != 0 {
	// 	panic("Unable to get question")
	// }

	return respJSON, nil
}

//TokenResp is the response body of a request for the token
type TokenResp struct {
	ResponseCode    int    `json:"response_code"`
	ResponseMessage string `json:"response_message"`
	Token           string `json:"token"`
}

func getToken() (token string) {
	resp, err := http.Get("https://opentdb.com/api_token.php?command=request")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	var tokenJSON TokenResp
	err = json.Unmarshal(respBody, &tokenJSON)
	if err != nil {
		return ""
	}

	return tokenJSON.Token
}

//ResetResp is the response body for our request to reset a token
type ResetResp struct {
	ResponseCode int    `json:"response_code"`
	Token        string `json:"token"`
}

func resetToken(token string) (err error) {
	resp, err := http.Get("https://opentdb.com/api_token.php?command=reset&token=" + token)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var resetJSON ResetResp
	err = json.Unmarshal(respBody, &resetJSON)
	if err != nil {
		return err
	}

	if resetJSON.ResponseCode != 0 {
		return fmt.Errorf("Resetting token failed, response code expected 0 got %d", resetJSON.ResponseCode)
	} else if resetJSON.Token != token {
		return fmt.Errorf("Resetting token failed, token mismatch")
	} else {
		return nil
	}
}

func shuffleify(CorrectAnswer string, IncorrectAnswers []string) (shuffled string) {
	allArr := append(IncorrectAnswers, CorrectAnswer)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(allArr), func(i, j int) { allArr[i], allArr[j] = allArr[j], allArr[i] })
	return strings.Join(allArr, ", ")
}

func getRoutine(qaPool chan<- TriviaElement) {
	token := ""
	for true {
		if token == "" {
			token = getToken()
		}

		qa, err := getTrivia(token)
		if err != nil {
			panic(err)
		}

		if qa.ResponseCode == 4 {
			//error what error
			go resetToken(token)
		} else if qa.ResponseCode == 3 {
			token = ""
		}

		for _, resElements := range qa.Results {
			que, err := base64.StdEncoding.DecodeString(resElements.Question)
			badErrHandler(err)
			ans, err := base64.StdEncoding.DecodeString(resElements.CorrectAnswer)
			badErrHandler(err)
			var incorrects []string
			for _, badAns := range resElements.IncorrectAnswers {
				badAnsDec, err := base64.StdEncoding.DecodeString(badAns)
				badErrHandler(err)
				incorrects = append(incorrects, string(badAnsDec))
			}
			qaPool <- TriviaElement{CorrectAnswer: string(ans), Question: string(que), PossibleAnswers: shuffleify(string(ans), incorrects)}
		}
	}
}

func askQuestion(qa TriviaElement) {
	ans, err := client.Prompt.Input("ans", fmt.Sprint(qa.Question, "\nAnswers might be one of ", qa.PossibleAnswers), ctoai.OptInputAllowEmpty(false))
	badErrHandler(err)
	if ans != qa.CorrectAnswer {
		printWrapper("BZZT! Wrong answer! Try again!")
		askQuestion(qa)
	} else {
		printWrapper("DING! You got it!")
	}
}

func main() {
	printLogo(client)
	
	qaPool := make(chan TriviaElement, 1)
	go getRoutine(qaPool)

	printWrapper("Time for a quiz kids!\nHere's your first question!")
	//0 = first time
	//1 = not first time, continue
	//-1 = quit
	stateInt := 0
	for stateInt >= 0 {
		if stateInt != 0 {
			printWrapper("Here's another question!")
		}

		askQuestion(<-qaPool)

		again, err := client.Prompt.Confirm("again", "Want another question?", ctoai.OptConfirmDefault(true))
		badErrHandler(err)
		if again {
			stateInt = 1
		} else {
			stateInt = -1
		}
	}
	printWrapper("Thank you for playing Ops Quiz!")
}
