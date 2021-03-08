package main
import(
	"html/template"
	"io/ioutil"
	"net/http"
	"log"
)
//A wiki consists of a series of interconnected pages, 
//each of which has a title and a body (the page content). 
type Page struct{
	Title string
	Body []byte
}
func (p *Page) Save() error{
	filename := p.Title + ".txt"
	return ioutil.WriteFile("./files/"+filename,p.Body,0600)
}
func LoadPage(title string) (*Page,error){
	filename := title+".txt"
	body,err:= ioutil.ReadFile("./files/"+filename)
	if err!=nil{
		return nil,err
	}
	return &Page{Title:title,Body:body},nil
}
func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/view/"):]
    p, err := LoadPage(title)
	if err != nil {
		http.Redirect(w,r,"/edit/"+title,http.StatusFound)
	}
    t, _ := template.ParseFiles("./views/view.html")
    t.Execute(w, p)
}
func saveHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title:title,Body:[]byte (body)}
	p.Save()
	http.Redirect(w,r,"/view/"+title,http.StatusFound)
}
func editHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/edit/"):]
	p,err := LoadPage(title)
	if err != nil {
		p = &Page{Title:title}
	}
	t,_:=template.ParseFiles("./views/edit.html")
	t.Execute(w,p)
}
func main(){
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./views/css/"))))
	http.HandleFunc("/view/",viewHandler)
	http.HandleFunc("/edit/",editHandler)
	http.HandleFunc("/save/",saveHandler)
	log.Fatal(http.ListenAndServe(":3000",nil))
}