   ============== PLANET ATTACK ===============

Conquer the galaxy!

A multiplayer simultanious-action turn-based game played out on a hex grid.  This is the base engine: galaxy creation and object interaction rules.  It needs:
    1: Cleaner Galaxy creation
    2: Rules/Creation Parameter Adjustments
    3: Snazzy javascript UI

Basic Rules:
    Planets have inhabitants, resources, and launchers.  Every turn each inhabitant will turn 1 resource into 1 launcher until resources are depleted.  

Players may order any number of available launchers fired toward another planet, where they become a Ship that travels uninterrupted in a straight linetoward its target a set distance each turn.  When a Ship hits its target (passing harmlessly through any Planets/Ships in the way), the Ship converts to inhabitants on the new Planet.

If a Ship hits a planet occupied by another faction (or neutral inhabitants), planet and cloud inhabitants cancel 1 for 1.  If the Ship has more inhabitants than the planet, ownership is transfered to the Ship player: otherwise (even if nobody is left alive), the planet's ownership is unchanged.

Each Ship is visible to any planet it passes near, but the only interaction possible right now is landing Ship on Planets.

Possible Rules Tweaks:
    Ships lose inhabitants over time
    Planets have population limits
    Ships do not launch on creation, must be given separate target orders
    Ships can be named
    Planet Resources slowly replenish
    "Racial" abilites available to different factions
