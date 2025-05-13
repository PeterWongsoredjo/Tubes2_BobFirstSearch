# Tubes2_BobFirstSearch# Little Alchemy 2 Recipe Search

## Algorithm Explanation

### **Depth-First Search (DFS)**
DFS is a search algorithm that explores a graph by following a single path until it reaches a base element, then backtracks and continues to explore other branches. DFS is implemented recursively, where each element is explored by following one of its components (children) to a certain depth. After reaching a base element, the algorithm backtracks to explore other components. This process allows DFS to deeply explore recipe paths, providing insight into all possible recipes that can be formed, even though these paths may be longer.

### **Breadth-First Search (BFS)**
BFS is a search algorithm that explores a graph level by level. Starting from the target element, BFS processes all components on the current level first, then moves deeper into the graph. This ensures that the recipe path found is the shortest because the algorithm prioritizes the exploration of elements closer to the start. Every element that has been visited is marked in the visited set to prevent redundancy. BFS ensures that the path to the target element can be found faster, especially when the shortest solution is required.

---

## Requirements

### **Program Requirements**

- **Go** version 1.17 or higher for backend development.
- **Node.js** version 14.x or higher for frontend development.
- **Docker** (optional) for containerized deployment.
- **Visual Studio Code** or any IDE with Go and React/Next.js support is recommended.

### **Installation**

1. **Clone the repository**:
   - Clone the repository and navigate to the project directory.
    ```bash
    git clone https://github.com/BobSwagg13/Tubes2_BobFirstSearch.git
    cd Tubes2_BobFirstSearch
    ```

2. **Frontend Installation**:
   - Navigate to the frontend directory.
   ```bash
   cd frontend
   ```
   - Install dependencies by running the appropriate command to download required packages.
    ```bash
    npm install
    ```
   - Start the frontend development server.
       ```bash
    npm run dev
    ```

3. **Backend Installation**:
   - Navigate to the backend directory.
    ```bash
    cd backend
    ```
   - Install Go dependencies by running the command to tidy the Go modules.
    ```bash
    go mod tidy
    ```

4. **Docker (Optional)**:
   - For containerization, use Docker Compose to build and start the containers.

---

## Commands to Run the Program

### **Backend**

1. To start the backend server, use the appropriate command below:

    ```bash
    go run ./main
    ```

2. To scrape the data from the Little Alchemy 2 website, use the following command:

    ```bash
    go run ./cmd
    ```

### **Frontend**
1. To start the frontend server locally, use the following command:

    ```bash
    npm run dev
    ```

## Author

**13523039 Peter Wongsoredjo**
*Email: 13523039@std.stei.itb.ac.id*
*GitHub: [PeterWongsoredjo](https://github.com/PeterWongsoredjo)*

**13523086 Bob Kunanda**
*Email: 13523086@std.stei.itb.ac.id*
*GitHub: [BobSwagg13](https://github.com/BobSwagg13)*

**13523109 Haegen Quinston**
*Email: 13523109@std.stei.itb.ac.id*
*GitHub: [haegenpro](https://github.com/haegenpro)*

*Mahasiswa Teknik Informatika, Institut Teknologi Bandung (ITB)*  
