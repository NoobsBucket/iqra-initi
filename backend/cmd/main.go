package main

func main(){

cng := config{
	address: ":8080",
	db: dbConfig{},
} 
api := application{
	config: cng,
}


}