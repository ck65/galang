package main

import (
	"crypto/rand"
	"fmt"
	"crypto/dsa"
	//"os"
	"io/ioutil"
	"crypto/ecdsa"
	"crypto/elliptic"
)

func main(){
	var t int
	var file string
	fmt.Println("input the file name:")
	fmt.Scanln(&file)
	fmt.Println("input the choice :1.ECDSA 2.DSA: ")
	fmt.Scanln(&t)
	if t==1{
		hash ,err:= ioutil.ReadFile(file)
		var c elliptic.Curve =elliptic.P224() 
		if err !=nil{
			fmt.Println(err)
		}else{
			fmt.Println(hash)
		priva,_ := ecdsa.GenerateKey(c,rand.Reader)
		s,r,err := ecdsa.Sign(rand.Reader, priva, hash)
		if err !=nil{
			fmt.Println(err)
		}else{
			if ecdsa.Verify(&priva.PublicKey,hash,s,r){
				fmt.Println("yes")
			}else{
				fmt.Println("no")
			}
			fmt.Println("the private key is:\n")
			fmt.Println(priva)
			}
		}
	}else if t==2{
		hash2 ,err := ioutil.ReadFile(file)
		if err !=nil{
			fmt.Println(err)
		}else{
			fmt.Println(hash2)
		priv := &dsa.PrivateKey{}
		dsa.GenerateParameters(&priv.Parameters, rand.Reader, dsa.L1024N160)
		dsa.GenerateKey(priv, rand.Reader)
		s,r,err := dsa.Sign(rand.Reader, priv, hash2)
		if err !=nil{
			fmt.Println(err)
		}else{
			if dsa.Verify(&priv.PublicKey,hash2,s,r){
				fmt.Println("yes")
			}else{
				fmt.Println("no")
			}
			fmt.Println("the private key is:\n")
			fmt.Println(priv)
			}
		}
	}
}