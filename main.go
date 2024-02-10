package main

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
)

type data struct {
	Time         string
	KeywordsJson string
	Keywords     string
	Extension    string
	Head         head
	Body         body
	Audio        string
	NameOfFolder string
	Source       string
}
type head struct {
	Title            string
	StoryDescription string
}
type body struct {
	Pages []pages
}
type pages struct {
	ImageLink       string
	FirstPage       bool
	Id              string
	ImageAlt        string
	TextDescription textt
}
type textt struct {
	Big   string
	Small string
}

func findTime() string {
	loc, e := time.LoadLocation("GMT")
	if e != nil {
		fmt.Println("Error in find time, determine location", e)
	}
	t := time.Now().In(loc).Format(time.DateTime)
	a := strings.Split(t, " ")
	t = a[0] + "T" + a[1] + "+00:00"
	return t
}

func input() data {
	var (
		title            string
		storyDescription string
	)
	var myMap = make(map[string]string)
	f, err := os.Open("data.cfg")
	if err != nil {
		fmt.Println("Error while Opening data.cfg", err.Error())
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		myMap[strings.Split(scanner.Text(), "=")[0]] = strings.Split(scanner.Text(), "=")[1]
	}
	title = myMap["Title"]
	storyDescription = myMap["StoryDescription"]
	var Head = head{
		Title:            title,
		StoryDescription: storyDescription,
	}

	imageAlt := "ImageAlt"
	big := "Big"
	small := "Small"

	imageNameArray, err := os.ReadDir("assets")
	if err != nil {
		fmt.Println("Error in reading image name", err.Error())
	}

	var pageArray []pages
	for i := 0; i < len(imageNameArray); i++ {
		firstPage := false
		if i == 0 {
			imageAlt += "1"
			big += "1"
			small += "1"
			firstPage = true
		} else {
			imageAlt = strings.Replace(imageAlt, string(imageAlt[len(imageAlt)-1]), strconv.Itoa(i+1), 1)
			big = strings.Replace(big, string(big[len(big)-1]), strconv.Itoa(i+1), 1)
			small = strings.Replace(small, string(small[len(small)-1]), strconv.Itoa(i+1), 1)
		}
		imageLink := "assets/" + (imageNameArray[i].Name())
		id := imageNameArray[i].Name()
		imageAltValue := myMap[imageAlt]
		smallValue := myMap[small]
		bigValue := myMap[big]
		texts := textt{
			Big:   bigValue,
			Small: smallValue,
		}

		page := pages{
			ImageLink:       imageLink,
			FirstPage:       firstPage,
			Id:              id,
			ImageAlt:        imageAltValue,
			TextDescription: texts,
		}
		pageArray = append(pageArray, page)

	}
	audioFileArray, err := os.ReadDir("audio")
	if err != nil {
		fmt.Println("Error in reading audio folder", err.Error())
	}
	var audioLink string
	if len(audioFileArray) == 0 {
		audioLink = ""
	} else {
		audioLink = "audio/" + audioFileArray[0].Name()
	}
	nameOfFolderVariable := myMap["NameOfFolderOrStory"]
	sourceVariable := myMap["Source"]
	keywords := myMap["Keywords"]
	keywordsjson := strings.Split(keywords, ",")
	var KeywordsJson string
	for a, v := range keywordsjson {
		if a != 0 {
			KeywordsJson += ","
		}
		KeywordsJson += "\"" + v + "\""
	}

	var Body = body{
		Pages: pageArray,
	}
	time := findTime()
	var Data = data{
		Time:         time,
		KeywordsJson: KeywordsJson,
		Keywords:     keywords,
		Head:         Head,
		Body:         Body,
		Audio:        audioLink,
		NameOfFolder: nameOfFolderVariable,
		Source:       sourceVariable,
	}
	return Data

}

func sitemapEditFormatting(folderName string, Time string) string {
	add := "\n<url>\n\t<loc>https://glamhub.co/webstories/" + folderName + "/</loc>\n\t<lastmod>" + Time + "</lastmod>\n\t<priority>0.80</priority>\n</url>\n</urlset>"
	return add

}

func sitemap(folderName string, Time string) {
	check, err := os.ReadDir("sitemap")
	if err != nil {
		fmt.Println("Error in reading sitemap folder ", err.Error())
		return
	}
	if len(check) == 0 {
		fmt.Println("\n\n\tNo sitemap detected! \n\n.")
		return
	}
	fmt.Println("\n\n\t** I have assumed that you have removed </urlset> from the sitemap **\n\n.")
	name := check[0].Name()
	f, err := os.OpenFile("sitemap/"+name, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("Error in opening sitemap file ", err.Error())
		return
	}
	defer f.Close()
	extra := sitemapEditFormatting(folderName, Time)
	_, err = f.WriteString(extra)
	if err != nil {
		fmt.Println("Error in appending to sitemap file ", err.Error())
		return
	}
}

func editIndexHtml(Data data) {
	filenamearray, err := os.ReadDir("html")
	if err != nil {
		fmt.Println("\n\n** Error in reading html directory **\n\n.", err)
		return
	}
	if len(filenamearray) == 0 {
		fmt.Println("\n\n** No Index.html File detected **\n\n.")
		return
	}
	fileName := filenamearray[0].Name()

	content, err := os.ReadFile("html/" + fileName)
	if err != nil {
		fmt.Println("\n\n** Error in reading content of index.html file **\n\n.", err)
		return
	}
	f, err := os.OpenFile("html/"+fileName, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("\n\n** Error opening file for writting **\n\n.", err)
		return
	}
	defer f.Close()
	content1 := strings.Split(string(content), "<!--$$-->")
	content2 := content1[0] + "\n<!--$$-->\n<a  class=\"shadow-2xl m-8 p-5 rounded-3xl flex justify-center \" href=\"https://glamhub.co/webstories/" + Data.NameOfFolder + "/\" target=\"_blank\"><span ><span class=\" flex justify-center \" ><img class=\"rounded-3xl md:w-9/12 w-11/12 aspect-h-16 aspect-w-9 \"src=\"" + Data.NameOfFolder + "/" + "assets/portrait." + Data.Extension + "\" alt=\"" + Data.Body.Pages[0].ImageAlt + "\"></span><span class=\" mt-4 flex justify-center \" ><h1 class=\" text-purple-700 md:text-xl font-bold \">" + Data.Head.Title + "</h1></span></span></a>\n" + content1[1]
	content = []byte(content2)

	_, err = f.Write(content)
	if err != nil {
		fmt.Println("\n\n** Error writing to the file index.html **\n\n.", err)
		return
	}
	fmt.Println("\n\n** Done writing file **\n\n.")

}
func crop() string {
	f, err := os.ReadDir("assets")
	if err != nil {
		fmt.Println("Error reading directory", err.Error())
	}
	imageName := f[0].Name()
	extension := strings.Split(imageName, ".")[1]
	img, err := imaging.Open("assets/" + imageName)
	if err != nil {
		fmt.Println("error opening image", err)
	}
	square := imaging.CropCenter(img, 640, 640)
	err = imaging.Save(square, "./assets/square."+extension)
	if err != nil {
		fmt.Println("error saving image", err)
	}
	square = imaging.CropCenter(img, 640, 853)
	err = imaging.Save(square, "./assets/portrait."+extension)
	if err != nil {
		fmt.Println("error saving image", err)
	}
	square = imaging.CropCenter(img, 853, 640)
	err = imaging.Save(square, "./assets/landscape."+extension)
	if err != nil {
		fmt.Println("error saving image", err)
	}
	fmt.Println("Done Image editing")
	return extension
}

func main() {
	Data := input()
	Data.Extension = crop()
	f, err := os.Create("index.html")
	if err != nil {
		fmt.Println("Encountered error creating index.html file", err.Error())
	}
	t, err := template.ParseFiles("templ.html")
	if err != nil {
		fmt.Println("Encountered error parsing templ.html file", err.Error())
	}
	err = t.Execute(f, Data)
	if err != nil {
		fmt.Println("Encountered error executing template", err.Error())
	}
	f.Close()

	//Updating Sitemap
	sitemap(Data.NameOfFolder, Data.Time)
	//Updating Index.html
	editIndexHtml(Data)

}
