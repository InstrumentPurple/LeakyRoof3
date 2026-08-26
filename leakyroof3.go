/* 2/27/2026 */
package main

import("fmt"
	"bufio"
	"time"
	"os"
	"strings"
	"strconv"
	"encoding/json"
	"math/rand"
)

var sc *bufio.Scanner

type StoreItem struct{
	Title string
	Desc string
	Price float64
	IsPM bool
	Stom FoodContent
}

type WorldState struct {
	ActiveMoca bool
	UnmocaDur time.Duration
	MocaDur time.Duration
	FoodAmt map[string]FoodContent
	Wallet float64
	StoreItems []StoreItem
	CharInv []StoreItem
	Tox FoodContent
	Energies map[Energy]float64
	Health int
}

type FoodContent struct{
	Silicate float64
	Iron float64
	Sludge float64
	Titanium float64
	Sugar float64
}


type attack struct {
	req Energy
	reqMuch int
	title string
	damage int
}


const (
	UNDEF = iota
	FIRE
	WATER
	PSY
	METAL
	ELECTRIC
	GREEN
)

type Energy int

type DidDie bool
type Command func(ws *WorldState)DidDie


func energyToString(instance Energy) string{
	switch instance {
		case FIRE:
			return "Fire"
		case WATER:
			return "Water"
		case PSY:
			return "Psy"
		case METAL:
			return "Metal"
		case ELECTRIC:
			return "Electric"
		case GREEN:
			return "Green";
		default:
			return "None"
	}
}

func startMenu(){
	fmt.Println("menu")
	fmt.Println("-------")
}

func scaleFloatPtr(dest *float64, src *float64, scalar float64){
	*dest = scalar * (*src)
}


func typeOut(str string, ws *WorldState){
	length := len(str)

	for i := 0; i < length;i++{
		fmt.Print(string(str[i]))
		if ws.ActiveMoca {
			time.Sleep(ws.MocaDur)
		} else {
			time.Sleep(ws.UnmocaDur)
		}
	}
	time.Sleep(time.Millisecond*130)
	fmt.Println("")
}


func gameExit(ws *WorldState){
	os.Exit(0)
}


func gameLoop(commands *map[string]Command, ws *WorldState){
	for fmt.Print(":"); sc.Scan(); fmt.Print(":"){
		cmd := sc.Text()
		cmd = strings.ToLower(cmd)
		fn, incmd := (*commands)[cmd]
		if incmd{
			dieState := fn(ws)
			if dieState{
				typeOut("you have died!", ws)
				gameExit(ws)
			}
		} else if cmd=="help" {
			for name,_ := range *commands{
				fmt.Println(name)
			}
		} else {
			typeOut("error command not found",ws)
		}
	}
}


func drinkMoca(ws *WorldState) DidDie {
	ws.ActiveMoca=!ws.ActiveMoca
	return false
}


func die(ws *WorldState) DidDie {
	return true
}


func feed(ws *WorldState) DidDie {
	startMenu()
	for foodName,content := range ws.FoodAmt{
		fmt.Println("##",foodName,"##")
		outStom(&content)
	}
	var infood bool
	var fa FoodContent
	for {
		fmt.Println()
		for _, item := range ws.CharInv{
			if item.IsPM{
				fmt.Println(item.Title)
			}
		}
		fmt.Print("Pokemon title to feed: ")
		sc.Scan()
		pktl := sc.Text()

		fmt.Print("food name:")
		sc.Scan()
		selected := sc.Text()
		fmt.Print("food ammount (float):")
		sc.Scan()
		amttxt := sc.Text()
		amt,_ := strconv.ParseFloat(amttxt,64)

		var selectedPM *StoreItem = nil
		for ind,item := range ws.CharInv{
			if item.IsPM && item.Title == pktl{
				selectedPM = &ws.CharInv[ind]
			}
		}

		if selectedPM == nil{
			fmt.Println("no such pokemon!")
			break;
		}

		fa, infood = ws.FoodAmt[selected]

		if !infood{
			fmt.Println("food not found")
		} else {
			/* add everything to the belly */
			scaleFloatPtr(&(selectedPM.Stom.Silicate),&(fa.Silicate), amt)
			scaleFloatPtr(&(selectedPM.Stom.Iron),&(fa.Iron), amt)
			scaleFloatPtr(&(selectedPM.Stom.Sludge),&(fa.Sludge), amt)
			scaleFloatPtr(&(selectedPM.Stom.Titanium),&(fa.Titanium), amt)
			scaleFloatPtr(&(selectedPM.Stom.Sugar),&(fa.Sugar), amt)
			break;
		}
	}

	return false
}


func outTox(t string){
	fmt.Println("Your pokemon is in " + t + " toxcity")
}


func checkTox(eat,tx float64, ty string)bool{
	if eat > tx{
		outTox(ty)
	}

	return eat > tx
}


func outStom(out *FoodContent){
	fmt.Println("silicate",out.Silicate)
	fmt.Println("iron",out.Iron)
	fmt.Println("sludge",out.Sludge)
	fmt.Println("titanium",out.Titanium)
	fmt.Println("sugar",out.Sugar)
}


func labs(ws *WorldState)DidDie{
	fmt.Println("nutrition: ")

	pmInd := make([]int,0)
	for ind, item := range ws.CharInv{
		if item.IsPM {
			pmInd = append(pmInd, ind)
		}
	}

	for _,ind := range pmInd{
		fmt.Println("##",ws.CharInv[ind].Title,"##")
		outStom(&ws.CharInv[ind].Stom)
		checkTox(ws.CharInv[ind].Stom.Silicate,ws.Tox.Silicate, "Silicate")
		checkTox(ws.CharInv[ind].Stom.Iron, ws.Tox.Iron, "Iron")
		checkTox(ws.CharInv[ind].Stom.Sludge, ws.Tox.Sludge, "Sludge")
		checkTox(ws.CharInv[ind].Stom.Titanium, ws.Tox.Titanium, "Titanium")
		checkTox(ws.CharInv[ind].Stom.Sugar, ws.Tox.Sugar, "Sugar")
	}
	return false
}


func pokestore(ws *WorldState)DidDie{
	startMenu()

	for num, item := range ws.StoreItems{
		fmt.Println()
		fmt.Println("id:",num)
		fmt.Println(item.Title)
		fmt.Println(item.Desc)
		fmt.Println(item.Price, "$")
	}
	fmt.Println()
	fmt.Println("you have ", ws.Wallet, " coin")

	fmt.Print("item to buy id = ")
	sc.Scan()
	id, _ := strconv.Atoi(sc.Text())
	if id >= 0 && id < len(ws.StoreItems){
		if ws.StoreItems[id].Price <= ws.Wallet{
			typeOut("Purchasing item: " + ws.StoreItems[id].Title, ws)
			ws.Wallet -= ws.StoreItems[id].Price
			ws.CharInv = append(ws.CharInv, ws.StoreItems[id])
		} else {
			typeOut("Not enough coin!", ws)
		}
	} else {
		fmt.Println("invalid selection")
	}

	return false
}


func save(ws *WorldState)DidDie{
	gameSaveData, _ := json.Marshal(*ws)
	timeStr := strconv.FormatInt(time.Now().Unix(), 10)
	os.Mkdir("./saves", 0777)
	newFile,_ := os.Create("./saves/" + timeStr)
	newFile.Write(gameSaveData)
	newFile.Close()
	return false
}


func loadSave(ws *WorldState)DidDie{
	got, _ := os.ReadDir("./saves/")
	for i, dir := range got{
		if !dir.IsDir(){
			fmt.Println(i,") ", dir.Name())
		}
	}

	fmt.Print("save id = ")
	sc.Scan()

	text := sc.Text()
	num, _ := strconv.Atoi(text)

	sel := got[num].Name()

	data, err := os.ReadFile("./saves/" + sel)
	if err != nil{
		fmt.Println("error with selection")
	}

	json.Unmarshal(data,ws)

	return false
}


func extractOne(elem *float64, value float64, extRatio float64) float64{
	prev := *elem
	*elem *= extRatio
	delta := prev - *elem
	return delta * value
}


func extractor(ws *WorldState)DidDie{
	silicateValue := 2.0
	silicateExtractionRatio := 0.8

	ironValue := 1.0
	ironExtractionRatio := 0.5

	titaniumValue := 0.2
	titaniumExtractionRatio := 0.92

	for _, item := range ws.CharInv{
		if item.IsPM{
			fmt.Println(item.Title)
		}
	}

	fmt.Print("selected pokemon = ")
	sc.Scan()
	pmt := sc.Text()

	var worth float64 = 0.0
	for ind,item := range ws.CharInv{
		if item.IsPM && item.Title == pmt{
			worth += extractOne(&(ws.CharInv[ind].Stom.Silicate), silicateValue, silicateExtractionRatio)
			worth += extractOne(&(ws.CharInv[ind].Stom.Iron), ironValue, ironExtractionRatio)
			worth += extractOne(&(ws.CharInv[ind].Stom.Titanium), titaniumValue, titaniumExtractionRatio)
		}
	}
	fmt.Println(worth, "worth of Silicates, Titanium, and Iron extracted")
	ws.Wallet += worth

	return false
}


func viewInv(ws *WorldState)DidDie{

	for _, item := range ws.CharInv{
		fmt.Println()
		fmt.Println(item.Title)
		fmt.Println(item.Desc)
	}

	return false
}


func gym(ws *WorldState)DidDie{
	typeOut("Select your pokemon to use for this battle: ",ws)
	for ind, item := range ws.CharInv{
		fmt.Print(ind)
		fmt.Print(" ) ")
		fmt.Println(item.Title)
	}
	sc.Scan()
	selText := sc.Text()
	sel, _ := strconv.Atoi(selText)


	yourAttacks := make([]attack,0)
	switch(ws.CharInv[sel].Title){
		case "Magnemite":
			yourAttacks = []attack{
				attack{title: "Lightning strike", damage: 50, reqMuch: 3, req: WATER },
				attack{title: "Divebomb", damage: 25, reqMuch: 0},
				attack{title: "Lightning ball", damage: 30, reqMuch: 1, req:ELECTRIC},
			}
		case "Umbreon":
			yourAttacks = []attack{
				attack{title: "Shadow strike", damage: 40, reqMuch: 2, req: PSY },
				attack{title: "Power Balls", damage: 35, reqMuch: 1, req: PSY},
				attack{title: "Kinetic spear", damage: 20, reqMuch: 0},
			}
		case "Miltank":
			yourAttacks = []attack{
				attack{title: "Milky abbyse", damage: 50, reqMuch: 2, req: WATER },
				attack{title: "Ice cream", damage: 37, reqMuch: 1, req:WATER},
				attack{title: "Hot Coco", damage: 30, reqMuch: 1, req:ELECTRIC},
			}
	}
	enemies := map[string]([]attack){
		"Charmander": []attack{
			attack{title: "Fire stream", damage: 15, reqMuch: 1, req:FIRE},
			attack{title: "Ember pilar", damage: 25, reqMuch: 1, req:FIRE},
			attack{title: "Smoke screen", damage: 30, reqMuch: 1, req:FIRE},
		},

		"Diglet": []attack{
			attack{title: "Earth bone", damage: 27},
			attack{title: "Stone grave", damage: 20},
			attack{title: "Rock tumble", damage: 15},
		},

		"Raichu": []attack{
			attack{title: "Thunder bolt", damage: 30},
			attack{title: "Hail storm lightning", damage: 40},
			attack{title: "Shock", damage: 15},
			attack{title: "Conduction", damage:13},
		},

	}

	typeOut("Select your opponent\n", ws)

	for name, _ := range enemies{
		fmt.Println(name)
	}


	var opponent []attack
	var okay bool = false

	var name string
	for !okay{
		sc.Scan()
		name = sc.Text()
		opponent, okay = enemies[name]
	}


	ws.Health = 100
	ohealth := 100

	anyTox := !checkTox(ws.CharInv[sel].Stom.Silicate,ws.Tox.Silicate, "Silicate") ||
		!checkTox(ws.CharInv[sel].Stom.Iron, ws.Tox.Iron, "Iron") ||
		!checkTox(ws.CharInv[sel].Stom.Sludge, ws.Tox.Sludge, "Sludge") ||
		!checkTox(ws.CharInv[sel].Stom.Titanium, ws.Tox.Titanium, "Titanium") ||
		!checkTox(ws.CharInv[sel].Stom.Sugar, ws.Tox.Sugar, "Sugar")

	if anyTox{
		fmt.Println("because your pokemon isn't healthy they will have less HP!")
		ws.Health = 70
	}

	for turn := 0; ws.Health > 0 && ohealth > 0; turn++{
		if turn % 2 == 0 {
			fmt.Println("Your turn")
			for i,s := range yourAttacks{
				fmt.Print(i, " ")
				typeOut(s.title, ws)
				fmt.Println("dmg: ", s.damage, "eng: ", s.reqMuch, energyToString(s.req))
			}

			fmt.Print(":")
			sc.Scan()
			opt,_ := strconv.Atoi(sc.Text())

			if opt < 0 || opt >= len(yourAttacks){
				typeOut("invalid Attack", ws)
				continue;
			}

			if ws.Energies[yourAttacks[opt].req] >= float64(yourAttacks[opt].reqMuch){
				ws.Energies[yourAttacks[opt].req] -= float64(yourAttacks[opt].reqMuch) // use energy
				ohealth -= yourAttacks[opt].damage
			} else {
				typeOut("Does not have enough energy to use that attack!", ws)
			}

		} else {
			randA := rand.Intn(len(opponent) * 1000) / 1000
			ws.Health -= opponent[randA].damage
			fmt.Println(name, " uses ")
			typeOut(opponent[randA].title, ws)
			fmt.Print("-")
			fmt.Println(opponent[randA].damage)
		}
	}

	if ws.Health > ohealth{
		typeOut("You won!", ws)
	} else {
		typeOut("You lost!",ws)
	}

	return false
}


func shower(ws *WorldState)DidDie{

	for ind, item := range ws.CharInv{
		if item.IsPM{
			ws.CharInv[ind].Stom.Sludge /=  2.0
		}
	}

	typeOut("showered", ws)
	return false
}


func newsPaper(ws *WorldState)DidDie{
	walstr := fmt.Sprintf("%v", ws.Wallet)
	typeOut( walstr + " Coin - 5", ws)
	if ws.Wallet - 5.0 >= 0.0{
		artc := []string{
			"New pokemon types invading Northern Forests. A recent evolution in Bulbasaur populations to the South quickly spread Northward last week causing powerlines to be destroyed and pushback from Scyther populations causing massive battles between armies of the opposing sides.",
		}

		sel := rand.Intn(len(artc))
		typeOut(artc[sel], ws)

	} else {
		typeOut("Not enough coin.", ws)
	}

	return false
}


func selectPokemon(ws *WorldState)int{
	for ind, item := range ws.CharInv{
		if item.IsPM{
			fmt.Print(ind)
			fmt.Print(" ) ")
			fmt.Println(item.Title)
		}
	}
	fmt.Print("Pokemon ID = ")
	sc.Scan()
	selText := sc.Text()
	sel, _ := strconv.Atoi(selText)

	return sel
}

func randomEnergy()Energy{
	return Energy(rand.Intn(int(GREEN)+1)) /* this is the last one in the const declaration */
}

func treadmill(ws *WorldState)DidDie{
	poke := selectPokemon(ws)

	var pokeStom *FoodContent = &( ws.CharInv[poke].Stom )
	pokeStom.Sugar *= 0.25

	ws.Energies[randomEnergy()] += float64(rand.Intn(6))
	ws.Energies[randomEnergy()] += float64(rand.Intn(6))

	return false
}


func checkEnergy(ws *WorldState)DidDie{
	for En, amt := range ws.Energies{
		fmt.Println(energyToString(En), amt)
	}

	return false
}



func main(){
	fmt.Println("Leaky roof 3: beyond the minerals")

	var ws WorldState

	ws.UnmocaDur,_ = time.ParseDuration("18ms")
	ws.MocaDur,_ = time.ParseDuration("10ms")
	ws.ActiveMoca = false
	ws.FoodAmt = map[string]FoodContent{
		"granite mix":{
			Silicate:1.2,
			Iron:4.0,
			Sludge:0.1,
			Titanium:20.0,
			Sugar: 4.0,
		},
		"water":{
			Silicate: 0.034,
			Iron: 0.009,
			Sludge:0.000934,
			Titanium: 0.001,
			Sugar:0.0,
		},
		"junk food":{
			Silicate: 0.064,
			Iron: 0.012,
			Sludge:0.000044,
			Titanium: 0.003,
			Sugar:2.04,
		},
		"rice":{
			Silicate: 0.08,
			Iron: 0.008,
			Sludge:0.00025,
			Titanium: 0.0,
			Sugar: 0.01,
		},
		"sushi":{
			Silicate: 0.0043,
			Iron: 0.002,
			Sludge:0.102,
			Titanium: 0.0010023,
			Sugar: 0.00852,
		},
	}
	ws.CharInv = append(ws.CharInv, StoreItem{Title:"Magnemite", Desc: "Pokemon. Electric type. Always your's.", IsPM: true})
	ws.StoreItems = []StoreItem{
		{Title:"Miltank",Desc:"Pokemon. This cow won't disapoint.", Price:2000.0, IsPM: true},
		{Title:"Umbreon",Desc:"Pokemon. This dark type is clever.", Price:2000.0, IsPM: true},
	}
	ws.Tox = FoodContent{Silicate:80.0, Iron: 60.0, Sludge:2.0, Titanium:90.0, Sugar: 500.0}
	ws.Energies = map[Energy]float64{}

	sc = bufio.NewScanner(os.Stdin)

	/* help not here. it's implemented in gameLoop */
	commands := map[string]Command{
		"moca":drinkMoca,
		"feed":feed,
		"labs":labs,
		"die":die,
		"pokestore":pokestore,
		"save":save,
		"load":loadSave,
		"extractor":extractor,
		"inv":viewInv,
		"gym":gym,
		"shower":shower,
		"news paper":newsPaper,
		"treadmill":treadmill,
		"check energy":checkEnergy,
	}

	gameLoop(&commands, &ws)
}
