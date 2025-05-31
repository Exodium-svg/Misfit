package Website

import (
	"MisFitsDiscord/Game"
	"MisFitsDiscord/Utils"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"net"
	"net/http"
	"strconv"
	"time"
)

var allowedAddresses = make(map[string]time.Time)

var env *Utils.Env
var session *discordgo.Session

func Start(_env *Utils.Env, _session *discordgo.Session) {
	env = _env
	session = _session
	fileServer := http.FileServer(http.Dir("Resource"))

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/env", getEnvHandler)
	http.HandleFunc("/save-env", saveEnvHandler)
	http.HandleFunc("/strings", getStringsHandler)
	http.HandleFunc("/save-strings", saveStringsHandler)
	http.HandleFunc("/save-ascensions", saveAscensionRoles)
	http.HandleFunc("/get-ascensions", getAscensionsHandler)
	http.HandleFunc("/get-roles", getServerRoles)

	http.Handle("/", fileServer)

	err := http.ListenAndServe(":8080", nil)

	if nil != err {
		panic(err)
	}
}

func isAllowed(request *http.Request) bool {
	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if nil != err {
		fmt.Printf("Error parsing remote address: %s\n", err)
		return false
	}

	allowedTime, exists := allowedAddresses[host]

	if !exists {
		return false
	}

	if allowedTime.Before(time.Now()) {
		delete(allowedAddresses, host)
		return false
	}

	return true
}

func loginHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := request.ParseForm()

	if nil != err {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte(err.Error()))
	}

	username := request.FormValue("username")
	password := request.FormValue("password")

	if Utils.EnvGet(env, "website.username", "misfits") != username {
		http.Redirect(writer, request, "/login.html", http.StatusForbidden)
		return
	} else if Utils.EnvGet(env, "website.password", "password") != password {
		http.Redirect(writer, request, "/login.html", http.StatusForbidden)
		return
	}

	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err == nil {
		allowedAddresses[host] = time.Now().Add(2 * time.Hour)
	}

	http.Redirect(writer, request, "/home-page.html", http.StatusFound)
}

func getStringsHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}

	if false == isAllowed(request) {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	jsonData, err := Utils.SerializeStrings()

	if nil != err {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = writer.Write(jsonData)

	if nil != err {
		fmt.Printf("Failed to write env to json err: %s\n", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getEnvHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if false == isAllowed(request) {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	jsonData, err := json.Marshal(Utils.EnvSerialize(env))

	if nil != err {
		fmt.Printf("Failed to encode env to json err: %s\n", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = writer.Write(jsonData)

	if nil != err {
		fmt.Printf("Failed to write env to json err: %s\n", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func saveEnvHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var entries []Utils.EnvEntry

	err := json.NewDecoder(request.Body).Decode(&entries)

	if nil != err {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, entry := range entries {
		switch entry.Type {
		case "string":
			Utils.EnvSet(env, entry.Key, entry.Value)
			break
		case "int":
			var val int
			val, err = strconv.Atoi(entry.Value)

			if nil != err {
				writer.WriteHeader(http.StatusBadRequest)
				return
			}

			Utils.EnvSet(env, entry.Key, val)
			break
		case "bool":
			var val bool
			val, err = strconv.ParseBool(entry.Value)
			if nil != err {
				writer.WriteHeader(http.StatusBadRequest)
			}
			Utils.EnvSet(env, entry.Key, val)

		case "float":
			var val float64
			val, err = strconv.ParseFloat(entry.Value, 64)
			if nil != err {
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			Utils.EnvSet(env, entry.Key, val)
			break
		}

		err = Utils.EnvSave(env)

		if nil != err {
			fmt.Printf("Failed to save env err: %s\n", err.Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

}

type StringKvp struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func saveStringsHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if false == isAllowed(request) {
		writer.WriteHeader(http.StatusForbidden)
		return
	}
	var entries []StringKvp

	err := json.NewDecoder(request.Body).Decode(&entries)

	if nil != err {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, entry := range entries {
		Utils.SetString(entry.Key, entry.Value)
	}

	Utils.SaveStrings()
}

func getAscensionsHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if false == isAllowed(request) {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	jsonData, err := json.Marshal(Game.GetAscensions())

	if nil != err {
		fmt.Printf("Failed to encode ascensions as json err: %s\n", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = writer.Write(jsonData)

	if nil != err {
		fmt.Printf("Failed to write ascensions to json err: %s\n", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func saveAscensionRoles(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if false == isAllowed(request) {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	var set []Game.AscensionSet
	err := json.NewDecoder(request.Body).Decode(&set)

	if nil != err {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	Game.SetAscensions(set)
}

func getServerRoles(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if false == isAllowed(request) {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	guildID := Utils.EnvGet(env, "discord.guild.id", "1370077540160372847")

	roles, err := session.GuildRoles(guildID)

	if nil != err {
		fmt.Printf("Failed to get guild roles err: %s\n", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(roles)

	if nil != err {
		fmt.Printf("Failed to encode roles as json err: %s\n", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = writer.Write(jsonData)
}
