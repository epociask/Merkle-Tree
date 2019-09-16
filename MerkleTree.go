package main

import(

	"crypto/sha512"
	"fmt"
//	"errors"
)

func hashFunction(str string) string{
	sha_512 := sha512.New()
	sha_512.Write([]byte(str))
	return string(sha_512.Sum(nil))
}



type Entry struct{
	Key string
	Value string
}

type MerkleTree struct{

	RootValue string
	Root *Node
	Count int
} 

type Node struct{
	Value string
	Left *Node
	Right *Node
	LeafValue *Entry
	Parent *Node

}

func(tree *MerkleTree) buildTree() {
	tree.Count = 0
	}
func (tree *MerkleTree) insert(key string, value string) string{
	//Instantiates LeafValue entry object for inserted node 
	entry := &Entry{key, value}
	//instantiate node to be inserted 
	insertNode := &Node{hashFunction(key+value), nil, nil, entry, nil}
	//Set temp equal to the Merkle Tree root for effective traversial 
	temp := tree.Root

	//Used to insert into empty tree or when count == 0 
	if tree.Count == 0{
		tree.Root = &Node{(key+value), nil, nil, entry, nil}

		tree.Count++
		return key

	//Base checks, sees if tree has two or three elements
	}else if tree.Count == 1 || tree.Count == 2{
		tree.Root = updatePath(temp, insertNode)
		tree.Count++
		return key

	}else if (tree.Count) % 4 == 0{
		//Checks if value is going to construct direct right node
		tree.Root = updatePath(temp, insertNode)
		tree.Count++
		return key
	}else{
	//All other condtions are dealt with iteration
	List := make([]*Node, int(tree.Count))
	X := true
	count := 1
	found := false

	List[0] = temp
	for{
		fmt.Println(count)
		if temp.Right.LeafValue != nil && found == false{
		found = true
		tempRoot := &Node{hashFunction(insertNode.Value+tree.Root.Value), temp, insertNode, nil, nil}
		temp.updateParent(tempRoot)
		insertNode.updateParent(tempRoot)
		List[count] = tempRoot
		X = false
		count--
		}

		if X == true {
		List[count] = temp
		temp = temp.Right

		count++


	}else if X == false && count >= 1{

		List[count] = updatePath(List[count-1].Left, List[count])
		count--

	}else if X == false && count == 0{
		tree.Root = updatePath(List[0].Left, List[1])
		tree.Count++
		return value
	}
	}

		return value

}
 
}

func (node *Node) updateParent(parent *Node){

	node.Parent = parent 
}
//@parm Left Node, Right Node 
//@returns parent hash node with pointers to left and righy children
func updatePath(leftNode, rightNode *Node) *Node{

	node := &Node{hashFunction(leftNode.Value + rightNode.Value), leftNode, rightNode, nil, nil}
	leftNode.updateParent(node)
	rightNode.updateParent(node)
	return node
	}

//@returns True ifleaf value is found 
func isLeaf(node *Node) bool{

	if node.LeafValue != nil {
		return true
	}else{
		return false
		}
	}

func findPathCount(node *Node) int {
	
	count := 0 
	for{
		if isLeaf(node) {
			return count

		}else{
			node = node.Left
			count++
		}
	}
	return count
	
}

func findNode(entry *Entry, node *Node) *Node{

	if node.Left != nil && node.Right != nil {
	return	findNode(entry, node.Right)
	return	findNode(entry, node.Left)

	} 

	if isLeaf(node) && (node.LeafValue.Value == entry.Value) && (node.LeafValue.Key == entry.Key) {

		return node

	}
	
	return nil
}

func generateMerklePath(entry Entry, node *Node) []string{
	et := &entry
	path := make([]string, findPathCount(node))



	temp := findNode(et, node)

	for i := 0; i > 0; i++{

		path[i] = node.Value

		if(temp.Parent == nil){
			break

		}else{

		temp = temp.Parent


		}
	}


	return path


	

}



func main(){

	var tree MerkleTree
	entry := Entry{"A", "$60"}
	tree.buildTree()
	tree.insert("A", "$60")
	tree.insert("B", "$300")
	tree.insert("C", "$45")	
	tree.insert("D", "$65")	
	tree.insert("E", "$UR DONE")
	path := generateMerklePath(entry, tree.Root)
	fmt.Println(path)
	fmt.Println(tree.Root.Value)
	fmt.Println(tree.Root.Left.Left.Left.Parent.Parent.Parent.Value)
//	fmt.Println(tree.Root.Right.LeafValue.Key)
//	fmt.Println(tree.Root.Right.Right.LeafValue.Key)
	//fmt.Println(tree.Root.Left.Right.Right.LeafValue.Key)
//		fmt.Println(tree.Root.Right.LeafValue.Key)
}