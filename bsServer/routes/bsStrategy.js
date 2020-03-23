/**
 * S19 CSCI 470 Web Science
 * 
 * battleship tile selection strategy
 * implements a random tile strategy
 *
 * Phillip J. Curtiss, Assistant Professor
 * pcurtiss@mtech.edu, 406-496-4807
 * Department of Computer Science, Montana Tech
 */
const tiles = []; // empty tiles array

// populate tiles array - row major
['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J'].forEach((row) => {
    ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9'].forEach((col) => {
        tiles.push(row + col);
    });
});

module.exports = {
    'getTile': function() {
        let tile = tiles[Math.floor(Math.random() * Math.floor(tiles.length))];
        let tileIdx = tiles.indexOf(tile);

        tiles.slice(tileIdx, 1);

        return (tile);
    }
};