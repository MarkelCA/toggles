package storage

import "github.com/markelca/toggles/flags"

type MemoryRepository struct {}

func NewMemoryRepository() flags.FlagRepository {
    return MemoryRepository{}
}

var flagsStorage []flags.Flag = make([]flags.Flag,0)

func(r MemoryRepository) Get(key string) (bool, error) {
    result,_:= r.List()
    for _,flag := range result {
        if flag.Name ==  key {
            return flag.Value,nil
        }
    }
    return false,nil
}

func (r MemoryRepository) List()([]flags.Flag, error) {
    return flagsStorage,nil
}

func (r MemoryRepository) Create(flag flags.Flag) error {
    result,err := r.Exists(flag.Name)

    if err != nil{
        return err
    }
    if result {
        return flags.FlagAlreadyExistsError
    }
    flagsStorage = append(flagsStorage,flag)
    return nil
}

func (r MemoryRepository) Exists(name string) (bool,error) {
    for _,currentFlag := range flagsStorage {
        if currentFlag.Name == name {
            return true,nil
        }
    }
    return false,nil
}

func (r MemoryRepository) Update(name string, value bool) error {
    for i := range flagsStorage {
        if flagsStorage[i].Name == name {
            flagsStorage[i].Value = value
            break
        }
    }
    return nil
}
