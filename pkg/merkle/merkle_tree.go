package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
)

type MerkleTree struct {
	Leaves []string
	Root   string
}

func NewMerkleTree() *MerkleTree {
	return &MerkleTree{}
}

func (mt *MerkleTree) AddLeaf(hash string) {
	mt.Leaves = append(mt.Leaves, hash)
	mt.Root = mt.computeRoot()
}

// Generates a Merkle proof for a given leaf hash
func (mt *MerkleTree) GenerateProof(hash string) ([]string, error) {
	index := -1
	for i, leaf := range mt.Leaves {
		if leaf == hash {
			index = i
			break
		}
	}
	if index == -1 {
		return nil, errors.New("hash not found in tree")
	}

	proof := []string{}
	layer := mt.Leaves
	for len(layer) > 1 {
		newLayer := []string{}
		for i := 0; i < len(layer); i += 2 {
			if i+1 == len(layer) {
				newLayer = append(newLayer, layer[i])
				if i == index {
					proof = append(proof, "")
				}
			} else {
				combined := combineHashes(layer[i], layer[i+1])
				newLayer = append(newLayer, combined)
				if i == index {
					proof = append(proof, layer[i+1])
				} else if i+1 == index {
					proof = append(proof, layer[i])
				}
			}
		}
		layer = newLayer
	}

	return proof, nil
}

func VerifyProof(hash string, proof []string, root string) bool {
	currentHash := hash
	for _, p := range proof {
		if p == "" {
			continue
		}
		currentHash = combineHashes(currentHash, p)
	}
	return currentHash == root
}

// Computes the Merkle root from the tree leaves
func (mt *MerkleTree) computeRoot() string {
	layer := mt.Leaves
	for len(layer) > 1 {
		newLayer := []string{}
		for i := 0; i < len(layer); i += 2 {
			if i+1 == len(layer) {
				newLayer = append(newLayer, layer[i])
			} else {
				newLayer = append(newLayer, combineHashes(layer[i], layer[i+1]))
			}
		}
		layer = newLayer
	}
	if len(layer) == 1 {
		return layer[0]
	}
	return ""
}

func HashFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func HashData(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func combineHashes(left, right string) string {
	data := left + right
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}