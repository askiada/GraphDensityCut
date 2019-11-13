# GraphDensityCut

Graph Clustering with Density-Cut
Junming Shao, Qinli Yang, Jinhu Liu and Stefan Kramer†

Understanding :

-   Build a Density-connected tree (DCT)
    Density Connectivity Map: DCT characterizes the
    density connectivity of vertices in graphs in a local
    fashion. It is intuitive that similar vertices are densely
    connected together, and vice versa
-   That tree is unique for each graph (see Theorem 1)
-   Each element of the DCT represents a component
-   We try to find the weakest edge in the DCT to create two partitions
-   We remove all the edges in the original graph which define these two partitions
-   The original graph know contains two partitions.
-   We repeat the same process from step 1 on each partition.
-   We stop when we are happy


## Documentation

https://godoc.org/github.com/askiada/GraphDensityCut