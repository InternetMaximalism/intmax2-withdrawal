package tree

import (
	"errors"
	"fmt"
	"intmax2-withdrawal/internal/hash/goldenposeidon"
)

type PoseidonHashOut = goldenposeidon.PoseidonHashOut

type PoseidonMerkleTree struct {
	height      uint8
	zeroHashes  []*PoseidonHashOut
	count       uint64
	siblings    []*PoseidonHashOut
	currentRoot PoseidonHashOut
}

// NewPoseidonMerkleTree creates new PoseidonMerkleTree by giving leaf nodes.
func NewPoseidonMerkleTree(
	height uint8,
	initialLeaves []*PoseidonHashOut,
	zeroHash *PoseidonHashOut,
) (mt *PoseidonMerkleTree, err error) {
	mt = &PoseidonMerkleTree{
		zeroHashes: generateZeroHashes(height, zeroHash),
		height:     height,
		count:      uint64(len(initialLeaves)),
	}

	mt.siblings, mt.currentRoot, err = mt.initSiblings(initialLeaves)
	if err != nil {
		return nil, errors.Join(ErrInitSISiblings, err)
	}

	return mt, nil
}

func buildIntermediate(leaves []*PoseidonHashOut) (nodes [][]*PoseidonHashOut, hashes []*PoseidonHashOut) {
	const (
		int0Key = 0
		int1Key = 1
		int2Key = 2
	)

	for i := int0Key; i < len(leaves); i += int2Key {
		var left, right = i, i + int1Key
		h := goldenposeidon.Compress(leaves[left], leaves[right])
		nodes = append(nodes, []*PoseidonHashOut{h, leaves[left], leaves[right]})
		hashes = append(hashes, h)
	}

	return nodes, hashes
}

// BuildMerkleRoot computes the root given the leaves of the tree
func (mt *PoseidonMerkleTree) BuildMerkleRoot(leaves []*PoseidonHashOut) (*PoseidonHashOut, error) {
	const (
		int0Key = 0
		int1Key = 1
		int2Key = 2
	)

	var (
		nodes [][][]*PoseidonHashOut
		ns    [][]*PoseidonHashOut
	)
	if len(leaves) == int0Key {
		leaves = append(leaves, mt.zeroHashes[int0Key])
	}

	for h := uint8(int0Key); h < mt.height; h++ {
		if len(leaves)%int2Key == int1Key {
			leaves = append(leaves, mt.zeroHashes[h])
		}
		ns, leaves = buildIntermediate(leaves)
		nodes = append(nodes, ns)
	}
	if len(ns) != int1Key {
		return nil, fmt.Errorf("%s: %+v", ErrBuildMerkleRootMoreThenOne.Error(), nodes)
	}

	return ns[int0Key][int0Key], nil
}

func generateZeroHashes(height uint8, zeroHash *PoseidonHashOut) []*PoseidonHashOut {
	const (
		int1Key = 1
	)
	var zeroHashes = []*PoseidonHashOut{
		new(PoseidonHashOut).Set(zeroHash),
	}
	// This generates a leaf = HashZero in position 0. In the rest of the positions that are equivalent to the ascending levels,
	// we set the hashes of the nodes. So all nodes from level i=5 will have the same value and same children nodes.
	for i := int1Key; i <= int(height); i++ {
		zeroHashes = append(zeroHashes, goldenposeidon.Compress(zeroHashes[i-int1Key], zeroHashes[i-int1Key]))
	}
	return zeroHashes
}

// ComputeMerkleProof computes the merkleProof and root given the leaves of the tree
func (mt *PoseidonMerkleTree) ComputeMerkleProof(index uint64, leaves []*PoseidonHashOut) (siblings []*PoseidonHashOut, root PoseidonHashOut, err error) {
	const (
		int0Key = 0
		int1Key = 1
		int2Key = 2
	)
	var ns [][]*PoseidonHashOut
	if len(leaves) == int0Key {
		leaves = append(leaves, mt.zeroHashes[int0Key])
	}
	proofIndex := index
	for h := uint8(int0Key); h < mt.height; h++ {
		if len(leaves)%int2Key == int1Key {
			leaves = append(leaves, mt.zeroHashes[h])
		}
		if proofIndex%int2Key == int1Key {
			// If it is odd
			siblings = append(siblings, leaves[proofIndex-int1Key])
		} else if len(leaves) > int1Key {
			if proofIndex >= uint64(len(leaves)) {
				siblings = append(siblings, leaves[proofIndex-int1Key])
			} else {
				siblings = append(siblings, leaves[proofIndex+int1Key])
			}
		}

		var (
			nsi    [][]*PoseidonHashOut
			hashes []*PoseidonHashOut
		)
		for i := int0Key; i < len(leaves); i += int2Key {
			var left, right = i, i + int1Key
			h := goldenposeidon.Compress(leaves[left], leaves[right])
			nsi = append(nsi, []*PoseidonHashOut{h, leaves[left], leaves[right]})
			hashes = append(hashes, h)
		}
		// Find the index of the leave in the next level of the tree.
		// Divide the index by 2 to find the position in the upper level
		proofIndex = uint64(float64(proofIndex) / int2Key)
		ns = nsi
		leaves = hashes
	}
	if len(ns) != int1Key {
		return nil, PoseidonHashOut{}, fmt.Errorf("%s: %+v", ErrBuildMerkleRootMoreThenOne, ns)
	}

	return siblings, *ns[int0Key][int0Key], nil
}

// AddLeaf adds new leaves to the tree and computes the new root
func (mt *PoseidonMerkleTree) AddLeaf(index uint64, leaf *PoseidonHashOut) (*PoseidonHashOut, error) {
	if index != mt.count {
		const msg = "mismatched leaf count: %d, expected: %d"
		return nil, fmt.Errorf(msg, index, mt.count)
	}
	cur := new(PoseidonHashOut).Set(leaf)
	isFilledSubTree := true

	const (
		int0Key = 0
		int1Key = 1
	)
	for h := uint8(int0Key); h < mt.height; h++ {
		if index&(int1Key<<h) > int0Key {
			child := cur
			parent := goldenposeidon.Compress(mt.siblings[h], child)
			cur = parent
		} else {
			if isFilledSubTree {
				// we will update the sibling when the sub tree is complete
				mt.siblings[h] = new(PoseidonHashOut).Set(cur)
				// we have a left child in this layer, it means the right child is empty so the sub tree is not completed
				isFilledSubTree = false
			}
			child := cur
			parent := goldenposeidon.Compress(child, mt.zeroHashes[h])
			cur = parent
			// the sibling of 0 bit should be the zero hash, since we are in the last node of the tree
		}
	}
	mt.currentRoot.Set(cur)
	mt.count++
	return cur, nil
}

// initSiblings returns the siblings of the node at the given index.
// it is used to initialize the siblings array in the beginning.
func (mt *PoseidonMerkleTree) initSiblings(initialLeaves []*PoseidonHashOut) (siblings []*PoseidonHashOut, root PoseidonHashOut, err error) {
	if mt.count != uint64(len(initialLeaves)) {
		return nil, PoseidonHashOut{}, ErrInitSiblingsFail
	}

	const (
		int0Key = 0
	)
	if mt.count == int0Key {
		for h := int0Key; h < int(mt.height); h++ {
			left := new(PoseidonHashOut)
			copy(left.Elements[:], mt.zeroHashes[h].Elements[:])
			siblings = append(siblings, left)
		}
		root, err := mt.BuildMerkleRoot(initialLeaves) // nolint:govet
		if err != nil {
			return nil, PoseidonHashOut{}, errors.Join(ErrCalculateInitMerkelRootFail, err)
		}
		return siblings, *root, nil
	}

	return mt.ComputeMerkleProof(mt.count, initialLeaves)
}

// GetCurrentRootCountAndSiblings returns the latest root, count and sibblings
func (mt *PoseidonMerkleTree) GetCurrentRootCountAndSiblings() (root PoseidonHashOut, count uint64, siblings []*PoseidonHashOut) {
	return mt.currentRoot, mt.count, mt.siblings
}

func (mt *PoseidonMerkleTree) CurrentRoot() PoseidonHashOut {
	return mt.currentRoot
}
