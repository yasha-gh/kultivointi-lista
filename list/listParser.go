package list

import (
	"context"
	"fmt"
	"kultivointi-lista/db"
	"kultivointi-lista/utils"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ListParser struct {
	ctx context.Context
	ListContent string
	List *List
}

func (l *ListParser) SetContext(ctx context.Context) {
	l.ctx = ctx
}

func (l *ListParser) LoadListFile() ([]*ListItem, error) {
	emptyList := []*ListItem{}
	filePath, err := runtime.OpenFileDialog(l.ctx, runtime.OpenDialogOptions{
		Title: "Avaa kultivointi lista",
		CanCreateDirectories: true,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Teksti tiedosto (*.txt)",
				Pattern: "*.txt",
			},
		},
	})
	if err != nil {
		return emptyList, err
	}
	listBytes, err := os.ReadFile(filePath)
	if err != nil {
		return emptyList, err
	}
	l.ListContent = string(listBytes)
	list, err := l.Parse()
	if err != nil {
		fmt.Println("Error parsing the list file", err)
	}
	dbConn, dbCtx, err := db.GetConn(l.ctx)
	if err != nil {
		return emptyList, err
	}
	defer dbConn.Close()
	tx, err := dbConn.BeginTx(dbCtx, nil)
	if err != nil {
		return emptyList, err
	}
	appCtx := l.ctx
	for _, listItem := range list {
		err := listItem.Save(appCtx, tx)
		if err != nil {
			fmt.Println("Error saving list item" , err)
		}
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println("Failed to commit list item saves")
	}

	return emptyList, nil
}

type ParseSerieRow struct {
	Titles []string `json:"titles"`
	LastSeenEpisodes []int `json:"lastSeenEpisodes"`
	Season int `json:"season"`
	BroadcastType string `json:"broadcastType"`
}

func (l *ListParser) Parse() ([]*ListItem, error) {
	lines := strings.Split(l.ListContent, "\r\n")
	orphanName := ""
	milExp := ""
	fullLines := []string{}
	log := utils.GetLogger()
	// Parse and combine off line titles
	for i, line := range lines {
		line := strings.TrimSpace(line)
		pattern := regexp.MustCompile(`\t+`)
		line = pattern.ReplaceAllString(line, "\t")
		lineParts := strings.Split(line, "\t")

		// Dunno why this is getting skipped
		if strings.Contains(line, "Shen Mu") {
			fullLines = append(fullLines, line)
		}

		if strings.Contains(line, "Zhanguo Qiannian") {
			milExp = line
			continue
		}


		if orphanName != "" {
			if orphanName == "Yao Shen Ji -" {
				orphanName = "Yao Shen Ji /"
			}
			// if orphanName == "With a Sword Domain," {
			// 	orphanName = "With a Sword Domain /"
			// }
			// Millenniums Of Warring States exception
			if milExp != "" {
				milParts := strings.Split(milExp, "\t")
				if len(milParts) == 3 {
					line = fmt.Sprintf("%s %s\t%v\t%v", milParts[0], orphanName, milParts[1], milParts[2])
				}
				// fmt.Println("mil exp line:", line)
				milExp = ""
				orphanName = ""
			} else {
				line = fmt.Sprintf("%s %s", orphanName, line)
				orphanName = ""
			}
		}

		if len(lineParts) < 2 {
			orphanName = line
		} else {
			// Ignore empty lines and header, remove extra tabs
			line = strings.TrimSpace(line)
			if line != "" && i > 2 {
				fullLines = append(fullLines, line)
			}
		}

	}

	fullLines = slices.DeleteFunc(fullLines, func(col string) bool {
		if col == "" {
			return true
		}
		return false
	})

	// testListBytes, _ := os.ReadFile("/home/defigo/Downloads/Kinuski animet ver 2-lines-only.txt")
	// testList := string(testListBytes)
	// testListSlice := strings.Split(testList, "\r\n")
	// testListTitles := []string{}
	// for _, tLine := range testListSlice {
	// 	tLine = strings.TrimSpace(tLine)
	// 	if tLine != "" {
	// 		parts := strings.Split(tLine, "\t")
	// 		if len(parts) > 0 {
	// 			testListTitles = append(testListTitles, strings.TrimSpace(parts[0]))
	// 		}
	// 	} else {
	// 		fmt.Println("Test file skipped:", tLine)
	// 	}
	// }
	// utils.PrettyPrint(fullLines)
	// fmt.Println("test titles count:", len(testListTitles))
	//
	// for _, title := range testListTitles {
	// 	foundTitle := false
	// 	for _, line := range fullLines {
	// 		parts := strings.Split(line, "\t")
	// 		fullListTitle := strings.TrimSpace(parts[0])
	// 		if title == fullListTitle {
	// 			foundTitle = true
	// 		}
	// 	}
	// 	if !foundTitle {
	// 		fmt.Println("Title missing", title)
	// 	}
	// }

	// list := List{
	// 	MainList: make([]*ListItem, 0),
	// }
	mainList := make([]*ListItem, 0)
	for _, line := range fullLines {
		parts := strings.Split(line,"\t")
		if len(parts) != 3 {
			if strings.Contains(line, "Three Swordsman Half Face") {
				titleParts := strings.Split(parts[0], " ")
				parts = []string{
					fmt.Sprintf("%s",strings.Join(titleParts[:len(titleParts[0])-1], " ")),
					"1",
					parts[1],
				}
			}
			if strings.Contains(line, "I Can Change The Timeline of Everything") {
				titleParts := strings.Split(parts[0], " ")
				parts = []string{
					fmt.Sprintf("%s", strings.Join(titleParts[:len(titleParts)-1], " ")),
					titleParts[len(titleParts)-1],
					parts[1],
				}
			}
			if strings.Contains(line, "Swallowed Star Movie:Blood Luo Continent") {
				titleParts := strings.Split(parts[0], " ")
				parts = []string{
					fmt.Sprintf("%s", strings.Join(titleParts[:len(titleParts)-1], " ")),
					titleParts[len(titleParts)-1],
					parts[1],
				}
			}
		}

		serieRow := ParseSerieRow{
			Titles: make([]string, 0),
			LastSeenEpisodes: make([]int, 0),
			Season: 0,
			BroadcastType: "ONA",
		}
		title := parts[0]
		if title == "Zhanshen Lianmeng , God of War Alliance" {
			title = "Zhanshen Lianmeng / God of War Alliance"
		}
		episode := parts[1]
		seasonStr := strings.ToLower(strings.TrimSpace(parts[2]))

		pthStart := strings.Index(title,"(")
		pthEnd := strings.Index(title,")")
		if pthStart != -1 && pthEnd != -1 {
			inPathensis := title[pthStart+1:pthEnd]
			if _, err := strconv.Atoi(inPathensis); err != nil {
				title = strings.Replace(title, "(", "/", 1)
				title = strings.Replace(title, ")", "", 1)
			}
		}

		titles := strings.Split(title, "/")
		for _, title := range titles {
			serieRow.Titles = append(serieRow.Titles, strings.TrimSpace(title))
		}

		episodesStr := strings.Split(episode, "/")
		for _, ep := range episodesStr {
			if ep == "7.12" {
				ep = "12"
			}
			ep, err := strconv.Atoi(strings.TrimSpace(ep))
			if err != nil {
				log.Error("Failed to convert episode number to int", "err", err)
			} else {
				serieRow.LastSeenEpisodes = append(serieRow.LastSeenEpisodes, ep)
			}
		}

		switch(seasonStr) {
		case "leffa":
			serieRow.Season = 1
			serieRow.BroadcastType = "movie"
		case "?", "-":
			serieRow.Season = 1
		default:
			season, err := strconv.Atoi(seasonStr)
			if err != nil {
				fmt.Println("Failed to convert season number to int", err)
			}
			serieRow.Season = season
		}

		seenOn := EpisodesSeen{}
		for _, ep := range serieRow.LastSeenEpisodes {
			seenOn = append(seenOn, &EpisodeSeen{
				EpisodesSeen: ep,
			})
		}
		primaryTitle := ""
		serieTitles := ItemTitles{}
		for i, title  := range serieRow.Titles {
			primary := true
			if i > 0 {
				primary = false
			}
			serieTitles = append(serieTitles, &ListItemTitle{
				Title: title,
				PrimaryTitle: primary,
			})
			if i == 0 {
				primaryTitle = title
			}
		}
		mainList = append(mainList, &ListItem{
			SeasonNum: serieRow.Season,
			Title: serieTitles[0],
			Titles: serieTitles,
			BroadcastType: serieRow.BroadcastType,
			EpisodesSeenOn: seenOn,
			Ongoing: true,
		})
		log.Info("Serie parsed", "title", primaryTitle)
	}

	log.Info("Import complete", "count", len(fullLines))
	return mainList, nil
}
