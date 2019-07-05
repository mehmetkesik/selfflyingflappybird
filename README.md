# self-flying flappy bird
self-flying flappy bird game with deep learning

# Usage
This project uses the **sdl2** graphics library. First you need to install the **sdl2** graphic library on your system and install the https://github.com/veandco/go-sdl2 library. later `go build` and run..

# Training
The artificial neural network receives the position, speed and height of the nearest pipe as the input. As the output gives the jump or jump. The entries of the artificial neural network are recorded. If the bird's life span is prolonged, the bird is rewarded and punished.
<br/><br/>
The artificial neural network is saved as the brain.json file in the main folder. If the brain.json file is deleted, the training starts from scratch.
