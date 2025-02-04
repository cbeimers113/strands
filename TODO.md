### v0.2.2 - Skybox & Better Highlighting:

- [ ] Fix GUI scaling in menus
- [ ] Notifications appear bottom up
- [ ] Horizon
- [ ] Pipeline improvements / build & release binaries / auto-tagging
- [x] Sky & stars
- [x] Real object highlighting instead of wireframe
- [x] Optimize chemical quantities tracker

### v0.3.0 - Plants Upgrade:

- Plants: they look silly as they grow; improve structure
  - Maybe be stages of trunk branching: each stage's branches are shorter, thinner and more numerous than the last
- Reproduction: make sure seeds and mutation is complete
- Photosynthesis: make them grow off of sunlight, carbon dioxide and water and produce Oxygen
  - Leaves produce more O2, take in more sunlight, CO2, H2O if higher leaf area in sun and color's distance from "ideal" green
- Nutrients/Toxins:
  - need levels of O2, CO2, H2O in the air
  - need to look into P, S, N2, other elements that could go in soil and water
- Produce fruit and/or thorns based on genetics
  - fruit will provice nourishment for creatures
  - thorns will deter creatures

### v0.4.0 - Creatures Upgrade:

- Creatures: randomly generated creatures that move around the map
- They will be able to eat the fruits of plants and the flesh of other creatures
  - will have genetics to determine how much nourishment is needed from fruit and other creatures
  - if a fruit is eaten, there is a chance the creature will poop out a seed with a copy of that plant's genetics (with potential mutations) a little while later
- Creatures will have bodies generated from simpler shapes and determined by genetics
- Strength and speed stats are also determined by genetics
  - higher strength and bigger size requires higher nutrition cost
  - higher speed also requires higher nutrition cost
- Basic neural network brain and state to determine actions

### v0.5.0 - Physics Upgrade:

- Water cycle: rain, evaporation, temperature simulation
  - wind
- Add woodiness to plants
- Fire and smoke

### v1.0.0 - First release version:

- Make map shape more irregular around edges and surround with water to make it like an island
- Saves server and save browser like TPT where users can share saved simulations
- Register account and login to upload saves to server

### Ideas for after v1.0.0

- More complex creature brains
- Better textures or models for plants and creatures
- Sample and save DNA of plants and animals
- DNA editor for plants and animals
- Re-investigate allowing saves to have different map dimensions
- Auto updates
