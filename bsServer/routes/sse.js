/**
 * S19 CSCI 470 Web Science
 * 
 * server sent events end points
 * send random tile guesses to the client over sse
 * triggering events based on distribution on tiles
 *
 * Phillip J. Curtiss, Assistant Professor
 * pcurtiss@mtech.edu, 406-496-4807
 * Department of Computer Science, Montana Tech
 */

const express = require('express'), // express module
    router = express.Router(); // define a route class

const tiles = []; // empty tiles array
var events = []; // empty events array

var ClientRes; // holds the Client Response object

// populate tiles array - row major
['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J'].forEach((row) => {
    ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9'].forEach((col) => {
        tiles.push(row + col);
    });
});

// populate events probability - same number as for tiles array
events = events.concat(Array.from({ length: 16 }).map(x => 'MISS'))
    .concat(Array.from({ length: 2 }).map(x => 'DESTROYER'))
    .concat(Array.from({ length: 16 }).map(x => 'MISS'))
    .concat(Array.from({ length: 3 }).map(x => 'SUBMARINE'))
    .concat(Array.from({ length: 16 }).map(x => 'MISS'))
    .concat(Array.from({ length: 3 }).map(x => 'CRUISER'))
    .concat(Array.from({ length: 16 }).map(x => 'MISS'))
    .concat(Array.from({ length: 4 }).map(x => 'BATTLESHIP'))
    .concat(Array.from({ length: 16 }).map(x => 'MISS'))
    .concat(Array.from({ length: 5 }).map(x => 'CARRIER'))
    .concat(Array.from({ length: 3 }).map(x => 'MISS'));

// send out message to the client
function sendSSEMSG(res, event, tile) {
    res.write(`event: ${event}\n`);
    res.write(`data: ${tile}\n\n`);
}

// End-point to establish a new stream
router.get('/stream', (req, res) => {
    // send the  correct content-type to establish the sse stream
    res.status(200).set({
        'content-Type': 'text/event-stream',
        'cache-control': 'no-cache',
        'connection': 'keep-alive'
    });

    ClientRes = res;

    // send a tile and its corresponding event every interval time
    let tileSelectionInterval = setInterval(() => {
        // compute a tile to send randomly
        let tile = tiles[Math.floor(Math.random() * Math.floor(tiles.length))];

        // obtain the index in the tiles array based on the tile selected
        let idx = tiles.indexOf(tile);

        // obtain the associated event from the tile selected
        let event = events[idx];

        // send the message using sse - event: and data: 
        sendSSEMSG(res, event, tile);

        // remove the selected tile from the tiles array
        tiles.splice(idx, 1);

        // remove the corresponding event from the events array
        events.splice(idx, 1);

        // stop selecting tiles once there are no more to select
        if (tiles.length == 0) {
            clearInterval(tileSelectionInterval);

            sendSSEMSG(res, 'END', '00');

            setTimeout(function() {
                res.end();
                ClientRes = {};
            }, 2000);
        }
    }, 500); // set interval to 5s (in ms)
}); // end-point stream


module.exports = {
    'router': router,
    'sseMessage': function(event, msg) {
        ClientRes.write(`event: ${event}\n`);
        ClientRes.write(`data: ${msg}\n\n`);
    }
};