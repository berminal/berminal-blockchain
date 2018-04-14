package console

import (
	"github.com/iost-official/prototype/tx/tx"
	"fmt"
	"github.com/iost-official/prototype/p2p"
	"strconv"
	"io/ioutil"
	"os"
	"github.com/iost-official/prototype/tx/min_framework"
	"github.com/iost-official/prototype/iostdb"
)

func Connect() Cmd {
	c := Cmd{
		Name:  "connect",
		Usage: `connect PORT - Connect to the network. Listen to PORT`,
	}
	c.Exec = func(host *Console, args []string) string {
		if len(args) != 1 {
			return "Invalid arguments!\n"
		}
		port, err := strconv.Atoi(args[0])
		if err != nil {
			return "Invalid arguments!\n"
		}

		dirname, _ := ioutil.TempDir(os.TempDir(), min_framework.DbFile)
		Db, err = iostdb.NewLDBDatabase(dirname, 0, 0)
		if err != nil{
			return "Can't open database"
		}

		Nn = p2p.NewNaiveNetwork()
		lis, err := Nn.Listen(uint16(port))
		if err != nil {
			return fmt.Sprint(err) + "\n"
		}

		Wg.Add(1)
		go func(<-chan p2p.Request, ) {
			defer Wg.Done()
			for {
				select{
				case message := <-lis:
					//fmt.Printf("\n%+v\n>", message)
					encodedBlock := message.Body
					block := transaction.DeserializeBlock(encodedBlock)
					err1 := Db.Put(block.Hash, encodedBlock)
					err2 := Db.Put([]byte("l"), block.Hash)
					if err1 != nil || err2 != nil {
						fmt.Printf("Write to database error! \nSync failed.\n>")
					}else{
						fmt.Printf("Sync successfully!\n>")
					}
				case <-Done:
					fmt.Printf("Port %d is done\n", port)
					return
				}
			}
		}(lis)
		return fmt.Sprintf("Connected with port %d successfully!\n", port)
	}
	return c
}
