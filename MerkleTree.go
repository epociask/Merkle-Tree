package main

import(

	"crypto/sha512"
	"fmt"
	"MerkleTree.go"
)

//Sha 512 has function
//@param string
//@returns hashString
func hashFunction(str string) string{

	sha_512 := sha512.New()
	sha_512.Write([]byte(str))
	return string(sha_512.Sum(nil))
}


//Struct to represent entry values as key/value pair
type Entry struct{

	Key string
	Value string
}

//Merkle Tree struct... holds rootValue String, root node, and count 
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
}

type MerklePath struct{

	Path []*Node
	Found bool
	Length int
}

//Used to instantiate merklePath
func buildPath() *MerklePath{
	
	return &MerklePath{nil, false, 0}
}
//Bread first algorithm 


//Function to instantiate Merkle Tree
func(tree *MerkleTree) BuildTree() {
	tree.Count = 0
	}


func (tree *MerkleTree) insert(key string, value string) string{
	//Instantiates LeafValue entry object for inserted node 
	entry := &Entry{key, value}
	//instantiate node to be inserted 
	insertNode := &Node{hashFunction(key+value), nil, nil, entry}
	//Set temp equal to the Merkle Tree root for effective traversial 
	temp := tree.Root

	//Used to insert into empty tree or when count == 0 
	if tree.Count == 0{
		tree.Root = &Node{(key+value), nil, nil, entry}

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
		if temp.Right.LeafValue != nil && found == false{
		found = true
		tempRoot := &Node{hashFunction(insertNode.Value+tree.Root.Value), temp, insertNode, nil}
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

//@parm Left Node, Right Node 
//@returns parent hash node with pointers to left and right children
func updatePath(Left, Right *Node) *Node{

	node := &Node{hashFunction(Left.Value + Right.Value), Left, Right, nil}
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

//@ param node & entry
//Checks to see if entry is within node
//If so we return true
func checkForEntry(node *Node, entry *Entry) bool{

	if (node.LeafValue.Value == entry.Value) && (node.LeafValue.Key == entry.Key) {

		return true
	}

	return false

}
//Checks to see if node has any children nodes.. if so we return true
func hasNoChildren(node *Node) bool {

	if node.Left == nil && node.Right == nil {

		return true
	}

		return false
}


//Recursive generate merkle path
//@param entry, node, list
//@recursively calls through tree by updating path values to a list, once specified path is found, store list as path
//If not found.. ie entire recursive stack executes.. return error string
func (merkle *MerklePath) generateMerklePath(entry *Entry, node *Node, list []*Node) string{

	if(merkle.Found == true){
		return "Found"
	}

	list = append(list, node, node.Left, node.Right)

	if (hasNoChildren(node) == true) && (checkForEntry(node, entry) == true){
		merkle.Found = true
		 merkle.Path = list[:len(list)-2]
		 merkle.Length = len(list) -2
		return "Found"
		
	}else if hasNoChildren(node) == true{
		

	}else{

	merkle.generateMerklePath(entry, node.Left, list[:len(list)-2])
	merkle.generateMerklePath(entry, node.Right, list[:len(list)-2])

}

		return "Error path not found"

}



//Function to find rightmost node in list
func findRight(root *Node) *Node{

	for {

		if(root.Right == nil){

			return root
		}

		root = root.Right
	}
}

//Takes two nodes, generates merkle paths for nodes and then swaps nodes by position
func (tree *MerkleTree) swapNodes(node1 *Node, node2 *Node){

	var passIn1, passIn2 []*Node
	path1 := buildPath()
	path1.generateMerklePath(node1.LeafValue, tree.Root, passIn1)
	path2 := buildPath()
	path2.generateMerklePath(node2.LeafValue, tree.Root, passIn2)
	path1.Length--
	path2.Length--
	temp1 := path1.Path[path1.Length]
	temp2 := path2.Path[path2.Length]


	if path1.Path[path1.Length-1].Right == temp1{

		path1.Path[path1.Length-1].Right = temp2

	}else if path1.Path[path1.Length-1].Left == temp1{

		path1.Path[path1.Length-1].Left = temp2

	}

	if path2.Path[path2.Length-1].Right == temp2 {

		path2.Path[path2.Length-1].Right = temp1

	}else if path2.Path[path2.Length-1].Left == temp2 {

		path2.Path[path2.Length-1].Left = temp1
	}

}


//List to hold collapsed tree
type List struct {
	L []*Node
}

//Function to instantiate list
func buildList() *List{

	return &List{nil}
}


//Method of list; finds all nodes in sequential order and appends them to the list effectively creating collapsed list
func (list *List) findNodes(root *Node) {

	if(root.LeafValue != nil){
		list.L = append(list.L, root)
		return
	}

	list.findNodes(root.Left)
	list.findNodes(root.Right)
	return
}

//Delete method of tree
func (tree *MerkleTree) delete(entry *Entry) {

	var passIn []*Node
	tempPath := buildPath()

	tempPath.generateMerklePath(entry, tree.Root, passIn)


	if(tempPath.Found == false){
		fmt.Println("ENTRY NOT FOUND")
		return
	}

	
	farRight := findRight(tree.Root)

	nodeToDelete := tempPath.Path[tempPath.Length-1]
	tree.swapNodes(nodeToDelete, farRight)	
	
	list := buildList()
	list.findNodes(tree.Root)

	var temp MerkleTree
	temp.BuildTree()

	//-1 to avoid last element within list 
	for i := 0; i <= len(list.L) -1; i++{

	temp.insert(list.L[i].LeafValue.Key, list.L[i].LeafValue.Value)
}

	
	tree.Root = temp.Root
	tree.RootValue = tree.Root.Value

	return
}	

//Checks given merkle path with generated one to ensure they are both equal; if not we return false
func (tree *MerkleTree) verifyMerklePath(entry *Entry, path *MerklePath) bool{

		var inPlace []*Node
		temp := buildPath()
		temp.generateMerklePath(entry, tree.Root, inPlace)

		if path == temp{

			return true

		}

		return false 
}

func main(){

	var tree MerkleTree
	entry := &Entry{"A", "$60"}
	tree.BuildTree()
	tree.insert("A", "$60")
	tree.insert("B", "$300")
	tree.insert("C", "$45")	
	tree.insert("D", "$65")	
	var temp []*Node
	path := buildPath()
	path.generateMerklePath(entry, tree.Root, temp)
	
	tree.delete(entry)


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
