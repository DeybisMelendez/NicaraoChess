# NicaraoChess

Nicarao Chess is a basic chess engine written in Go by Deybis Melendez.

[Lichess Account](https://lichess.org/@/Nicarao_Chess) (2100 elo blitz).

## Features

### Communication 

- UCI Protocol

### Board

- github.com/dylhunn/dragontoothmg

### Search

- Principal Variation Search
- Iterative Deepening
- Transposition Table

- Move Ordering:
    - Hash Move
    - PV Move
    - Killer Heuristic
    - History Heuristic
    - MVV-LVA
    - SEE

- Selectivity:
    - Check Extension
    - Late Move Reduction (LMR)
    - Futility Pruning
    - Quiescence
    - Mate Distance Pruning

### Evaluation

- Material
- Piece Square Table
- Evaluation of pieces
- Tapered eval
- Basic Draw Evaluation

### TODO

- Aspiration Window
- Draw Evaluation
- Mop-up Evaluation
- Internal Iterative Deepening
- Countermove Heuristic