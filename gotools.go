// ---------------------------------------------------------------------- go toolbox - mysql rest api (mwx'2024)
package gotools

import (
  "fmt"
  "os"
  "io"
  "net"
  "strconv"
  "crypto/rand"
 	"math/big"
  "database/sql"
  "github.com/spf13/viper"
  "github.com/AlecAivazis/survey/v2"
  "github.com/AlecAivazis/survey/v2/terminal"
  "flag"
  "net/http"
  "github.com/fatih/color"
  "encoding/base64"
  "github.com/eknkc/basex"
)

var Cr func(...interface {}) string = color.New(color.FgRed ).SprintFunc()
var Cg func(...interface {}) string = color.New(color.FgGreen).SprintFunc()
var Cy func(...interface {}) string = color.New(color.FgYellow).SprintFunc()
var Cb func(...interface {}) string = color.New(color.FgBlue).SprintFunc()
var Cm func(...interface {}) string = color.New(color.FgMagenta).SprintFunc()
var Cc func(...interface {}) string = color.New(color.FgCyan).SprintFunc()
var Cw func(...interface {}) string = color.New(color.FgWhite).SprintFunc()

var Crb func(...interface {}) string = color.New(color.Bold,color.FgRed ).SprintFunc()
var Cgb func(...interface {}) string = color.New(color.Bold,color.FgGreen).SprintFunc()
var Cyb func(...interface {}) string = color.New(color.Bold,color.FgYellow).SprintFunc()
var Cbb func(...interface {}) string = color.New(color.Bold,color.FgBlue).SprintFunc()
var Cmb func(...interface {}) string = color.New(color.Bold,color.FgMagenta).SprintFunc()
var Ccb func(...interface {}) string = color.New(color.Bold,color.FgCyan).SprintFunc()
var Cwb func(...interface {}) string = color.New(color.Bold,color.FgWhite).SprintFunc()

var CHRS = "VW9IdGJ6eXh1T25DRHdrc2M5MlhOQVNQcEJFWnJhWVY2ZEowaFJLdmoxNUdxVDRJZkZpTTdRZW0zTFc4Z2w="
var LR = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func enc(str string) string {
  enc,_ := basex.NewEncoding(db64(CHRS))
  return enc.Encode([]byte(rndstr(2)+str))
}

func dec(str string) string {
  enc,_ := basex.NewEncoding(db64(CHRS))
  b,_:=enc.Decode(str)
  return string(b[2:])
}

func db64(txt string) string {
  d,_ := base64.StdEncoding.DecodeString(txt);
  return string(d)
}

func rndstr(n int) string {
  b := make([]rune, n)
  for i := range b {
    n,_ := rand.Int(rand.Reader, big.NewInt(int64(len(LR))))
    b[i] = LR[n.Int64()]
  }
  return string(b)
  return""
}

func input(msg string, def string) string { // -------------------------------- AlecAivazis/survey: input string
  tmp:="";
  err:=survey.AskOne(&survey.Input{Message:msg,Default:def},&tmp)
  if err != nil {	
    if err==terminal.InterruptErr { 
      P(Crb("Interrupted."))
      os.Exit(0) 
    }
  }
  return tmp
}

func inputpw(msg string) string { // ---------------------------------------- AlecAivazis/survey: input password
  tmp:="";
  err:=survey.AskOne(&survey.Password{Message:msg}, &tmp)
  if err != nil {	
    if err==terminal.InterruptErr { 
      P(Crb("Interrupted."))
      os.Exit(0) 
    }
  }
  return tmp
}


func yesno(msg string, def bool) bool { // ------------------------------------------ AlecAivazis/survey: yes/no
  var err error
  tmp:=""
  if (def) {
    err=survey.AskOne(&survey.Select{ Message:msg, Options: []string{"Yes", "No"},},&tmp)
  } else {
    err=survey.AskOne(&survey.Select{ Message:msg, Options: []string{"No", "Yes"},},&tmp)
  }
  
  if err != nil {	
    if err==terminal.InterruptErr { 
      P(Crb("Interrupted."))
      os.Exit(0) 
    }
  }

  if (tmp=="Yes") { return true } else { return false }
}

func getid(n int) (string) { // ------------------------------------------------------- get base62 random string
  const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
  ret := make([]byte, n)
  for i := 0; i < n; i++ {
    num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
    ret[i] = letters[num.Int64()]
  }
  return string(ret)
}

func mytable(rows *sql.Rows) []map[string]interface{} { // ----------------------------- load mysql result table
  defer rows.Close()
  columns, _ := rows.Columns()
  count := len(columns)
  tableData := make([]map[string]interface{}, 0)
  values := make([]interface{}, count)
  valuePtrs := make([]interface{}, count)
  for rows.Next() {
    for i := 0; i < count; i++ { valuePtrs[i] = &values[i] }
    rows.Scan(valuePtrs...)
    entry := make(map[string]interface{})
    for i, col := range columns {
      var v interface{}
      val := values[i]
      b, ok := val.([]byte)
      if ok { v = string(b) } else { v = val }
      entry[col] = v
    }
    tableData = append(tableData, entry)
  }
  return(tableData)
}
 
func checkip(network string, ip string) bool { // ---------- check if ip is in range (cidr address or single ip)
  if net.ParseIP(ip) == nil { return false }
  _, subnet, err := net.ParseCIDR(network)
  if (err==nil) {
    if subnet.Contains(net.ParseIP(ip)) {
      return true
    }
  } else {
    if (network==ip) { return true }
  }
  return false
}

func isflagpassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func body(r *http.Response) (string) {
  body, err := io.ReadAll(r.Body)
  if (err==nil) {
    return string(body)
  }
  return ""
}

func P(a ...any) (n int, err error) { return fmt.Fprintln(os.Stdout, a...) }
func PN(a ...any) (n int, err error) { return fmt.Fprint(os.Stdout, a...) }
func PF(format string, a ...any) (n int, err error) { return fmt.Fprintf(os.Stdout, format, a...) }
func SF(format string, a ...any) string { return fmt.Sprintf(format, a...) }

func atoi(s string) int {
  i,err:=strconv.Atoi(s)
  if err != nil { return 0 }
  return i
}

func itoa(i int) string { return strconv.Itoa(i) }

func VX(key string) bool { return viper.IsSet(key) }
func VS(key string) string { return viper.GetString(key) }
func VSS(key string) []string { return viper.GetStringSlice(key) }
func VB(key string) bool { return viper.GetBool(key) }
func VI(key string) int { return viper.GetInt(key) }

func PE(msg string) (n int, err error) { return fmt.Fprintf(os.Stdout, "%s: %s\n",Crb("ERROR"),Cwb(msg)) }
func PO(msg string) (n int, err error) { return fmt.Fprintf(os.Stdout, "%s: %s\n",Cgb("OK"),Cwb(msg)) }


// --------------------------------------------------------------------------------------------------------- END
