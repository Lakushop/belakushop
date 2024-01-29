package lakushop

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/lakushop/belakushop/model"
	"github.com/lakushop/belakushop/module"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/argon2"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

var db = module.MongoConnect("MONGOSTRING", "lsdb")

func TestGetUserFromEmail(t *testing.T) {
	email := "admin@gmail.com"
	hasil, err := module.GetUserFromEmail(email, db)
	if err != nil {
		t.Errorf("Error TestGetUserFromEmail: %v", err)
	} else {
		fmt.Println(hasil)
	}
}

func TestInsertOneproduct(t *testing.T) {
	var doc model.Product
	doc.NamaProduct = "Baju anak motif spongebob"
	doc.Deskripsi = "baju anak lembut dengan motif menarik"
	doc.Kategori = "Baju anak"
	doc.Harga = "RP 120.0000"
	if doc.NamaProduct == "" || doc.Deskripsi == "" || doc.Kategori == "" || doc.Harga == "" {
		t.Errorf("mohon untuk melengkapi data")
	} else {
		insertedID, err := module.InsertOneDoc(db, "product", doc)
		if err != nil {
			t.Errorf("Error inserting document: %v", err)
			fmt.Println("Data tidak berhasil disimpan")
		} else {
			fmt.Println("Data berhasil disimpan dengan id :", insertedID.Hex())
		}
	}
}

type Userr struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Email string             `bson:"email,omitempty" json:"email,omitempty"`
	Role  string             `bson:"role,omitempty" json:"role,omitempty"`
}

func TestGetAllDoc(t *testing.T) {
	hasil := module.GetAllDocs(db, "user", []Userr{})
	fmt.Println(hasil)
}

func TestInsertUser(t *testing.T) {
	var doc model.User
	doc.Email = "admin@gmail.com"
	password := "admin123"
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		t.Errorf("kesalahan server : salt")
	} else {
		hashedPassword := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
		user := bson.M{
			"email":    doc.Email,
			"password": hex.EncodeToString(hashedPassword),
			"salt":     hex.EncodeToString(salt),
			"role":     "admin",
		}
		_, err = module.InsertOneDoc(db, "user", user)
		if err != nil {
			t.Errorf("gagal insert")
		} else {
			fmt.Println("berhasil insert")
		}
	}
}

func TestGetUserByAdmin(t *testing.T) {
	id := "65b7a3b39c00a28eabf478e8"
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		t.Errorf("Error converting id to objectID: %v", err)
	}
	data, err := module.GetUserFromID(idparam, db)
	if err != nil {
		t.Errorf("Error getting document: %v", err)
	} else {
		if data.Role == "pengguna" {
			datapengguna, err := module.GetPenggunaFromAkun(data.ID, db)
			if err != nil {
				t.Errorf("Error getting document: %v", err)
			} else {
				datapengguna.Akun = data
				fmt.Println(datapengguna)
			}
		}
		if data.Role == "seller" {
			dataseller, err := module.GetSellerFromAkun(data.ID, db)
			if err != nil {
				t.Errorf("Error getting document: %v", err)
			} else {
				dataseller.Akun = data
				fmt.Println(dataseller)
			}
		}
	}
}

func TestSignUpPengguna(t *testing.T) {
	var doc model.Pengguna
	doc.NamaLengkap = "megah"
	doc.TanggalLahir = "21/12/2000"
	doc.JenisKelamin = "laki-laki"
	doc.NomorHP = "080000000000"
	doc.Alamat = "Bandung"
	doc.Akun.Email = "megah1@gmail.com"
	doc.Akun.Password = "admin123"
	err := module.SignUpPengguna(db, doc)
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	} else {
		fmt.Println("Data berhasil disimpan dengan nama :", doc.NamaLengkap)
	}
}

func TestLogIn(t *testing.T) {
	var doc model.User
	doc.Email = "admin@gmail.com"
	doc.Password = "admin123"
	user, err := module.LogIn(db, doc)
	if err != nil {
		t.Errorf("Error getting document: %v", err)
	} else {
		fmt.Println("Selamat datang user:", user)
	}
}

func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := module.GenerateKey()
	fmt.Println("ini private key :", privateKey)
	fmt.Println("ini public key :", publicKey)
	id := "6569a026a943657839880665"
	objectId, err := primitive.ObjectIDFromHex(id)
	role := "pengguna"
	if err != nil {
		t.Fatalf("error converting id to objectID: %v", err)
	}
	hasil, err := module.Encode(objectId, role, privateKey)
	fmt.Println("ini hasil :", hasil, err)
}

func TestUpdatePengguna(t *testing.T) {
	var doc model.Pengguna
	id := "65b7a3b39c00a28eabf478e8"
	objectId, _ := primitive.ObjectIDFromHex(id)
	id2 := "65b7a3b39c00a28eabf478e8"
	userid, _ := primitive.ObjectIDFromHex(id2)
	doc.NamaLengkap = "Admin"
	doc.TanggalLahir = "21/12/2000"
	doc.JenisKelamin = "laki-laki"
	doc.NomorHP = "0800000000"
	doc.Alamat = "Bandung"
	if doc.NamaLengkap == "" || doc.TanggalLahir == "" || doc.JenisKelamin == "" || doc.NomorHP == "" || doc.Alamat == "" {
		t.Errorf("mohon untuk melengkapi data")
	} else {
		err := module.UpdatePengguna(objectId, userid, db, doc)
		if err != nil {
			t.Errorf("Error inserting document: %v", err)
			fmt.Println("Data tidak berhasil diupdate")
		} else {
			fmt.Println("Data berhasil diupdate")
		}
	}
}

func TestWatoken(t *testing.T) {
	body, err := module.Decode("fca3dbba6c382d6e937d33837f7428c1211e01a9928cbbbc0b86bb8351c02407", " v4.public.eyJleHAiOiIyMDIzLTEyLTAxVDE4OjU4OjE1KzA4OjAwIiwiaWF0IjoiMjAyMy0xMi0wMVQxNjo1ODoxNSswODowMCIsImlkIjoiNjU1YzNiOWExZDY1MjRmMmYxMjAwZmM2IiwibmJmIjoiMjAyMy0xMi0wMVQxNjo1ODoxNSswODowMCIsInJvbGUiOiJwZW5nZ3VuYSJ9GIKgKcp8gj4lzPH_NFvpx3GR2kBZ2qsDquYMKQdQ1PFpvHKlDy-FeO1umIGCaMuYyACP5jd-Y0at1WCOrsNRCA")
	fmt.Println("isi : ", body, err)
}

func TestInsertOneOrder(t *testing.T) {
	var doc model.Orderproduct
	doc.NamaProduct = "Event coldplay"
	doc.Quantity = "1"
	doc.TotalCost = "Rp 1000.000"
	doc.Status = "Pending"
	if doc.Quantity == "" || doc.TotalCost == "" || doc.Status == "" {
		t.Errorf("mohon untuk melengkapi data")
	} else {
		insertedID, err := module.InsertOneDoc(db, "order", doc)
		if err != nil {
			t.Errorf("Error inserting document: %v", err)
			fmt.Println("Data tidak berhasil disimpan")
		} else {
			fmt.Println("Data berhasil disimpan dengan id :", insertedID.Hex())
		}
	}
}

func TestUpdateProduct(t *testing.T) {
	conn := module.MongoConnect("MONGOSTRING", "lsdb")
	payload, err := module.Decode("", "v4.public.eyJleHAiOiIyMDIzLTEyLTA0VDE2OjM4OjUzKzA3OjAwIiwiaWF0IjoiMjAyMy0xMi0wNFQxNDozODo1MyswNzowMCIsImlkIjoiNjU2OWEwMjZhOTQzNjU3ODM5ODgwNjY1IiwibmJmIjoiMjAyMy0xMi0wNFQxNDozODo1MyswNzowMCIsInJvbGUiOiJwZW5nZ3VuYSJ97W3y3P-Q0NPzHRef8UNVz1eDQ-Ucx3_vDm23Gb6XicxGm4B0LTAYcA8Q7v1Nl_MXJIb9XATP70-6URg8zYVtCA")
	if err != nil {
		t.Errorf("Error decode token: %v", err)
	}
	if payload.Role != "admin" {
		t.Errorf("Error role: %v", err)
	}
	var datatiket model.Product
	datatiket.NamaProduct = "Event Coldplay 3 surabaya"
	datatiket.Deskripsi = "Terminal bus surabaya "
	datatiket.Kategori = "jam jemputan 13:00"
	datatiket.Harga = "Rp 100.000"
	id := "6569a53d783c6970079a560b"
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		t.Fatalf("error converting id to objectID: %v", err)
	}
	err = module.Updateproduct(objectId, payload.Id, conn, datatiket)
	if err != nil {
		t.Errorf("Error update : %v", err)
	} else {
		fmt.Println("Success!!!")
	}
}

func TestDeleteProduct(t *testing.T) {
	conn := module.MongoConnect("MONGOSTRING", "lsdb")
	payload, err := module.Decode("fca3dbba6c382d6e937d33837f7428c1211e01a9928cbbbc0b86bb8351c02407", "v4.public.eyJleHAiOiIyMDIzLTEyLTAxVDE4OjU4OjE1KzA4OjAwIiwiaWF0IjoiMjAyMy0xMi0wMVQxNjo1ODoxNSswODowMCIsImlkIjoiNjU1YzNiOWExZDY1MjRmMmYxMjAwZmM2IiwibmJmIjoiMjAyMy0xMi0wMVQxNjo1ODoxNSswODowMCIsInJvbGUiOiJwZW5nZ3VuYSJ9GIKgKcp8gj4lzPH_NFvpx3GR2kBZ2qsDquYMKQdQ1PFpvHKlDy-FeO1umIGCaMuYyACP5jd-Y0at1WCOrsNRCA")
	if err != nil {
		t.Errorf("Error decode token: %v", err)
	}
	id := "6569a53d783c6970079a560b"
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		t.Fatalf("error converting id to objectID: %v", err)
	}
	err = module.DeleteProduct(objectId, payload.Id, conn)
	if err != nil {
		t.Errorf("Error delete : %v", err)
	} else {
		fmt.Println("Success!!!")
	}
}

func TestGetAllProduct(t *testing.T) {
	conn := module.MongoConnect("MONGOSTRING", "lsdb")
	data, err := module.GetAllProduct(conn)
	if err != nil {
		t.Errorf("Error get all : %v", err)
	} else {
		fmt.Println(data)
	}
}

func TestGetProductFromID(t *testing.T) {
	conn := module.MongoConnect("MONGOSTRING", "lsdb")
	id := "6569a025a943657839880661"
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		t.Fatalf("error converting id to objectID: %v", err)
	}
	product, err := module.GetProductFromID(objectId, conn)
	if err != nil {
		t.Errorf("Error get Product : %v", err)
	} else {
		fmt.Println(product)
	}
}

func TestReturnStruct(t *testing.T) {
	id := "11b98454e034f3045021a8aa8eb84280"
	objectId, _ := primitive.ObjectIDFromHex(id)
	user, _ := module.GetUserFromID(objectId, db)
	data := model.User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	hasil := module.GCFReturnStruct(data)
	fmt.Println(hasil)
}
