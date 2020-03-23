/**
 * S19 CSCI 470 Web Science
 * 
 * battleship protocol 
 * Game Playing protocol between two node servers
 * 
 * DEFINITIONS:
 *  Opponent: The remote Node Server 
 *  Client: The web-browser client associated with this Node Server
 *
 * Phillip J. Curtiss, Assistant Professor
 * pcurtiss@mtech.edu, 406-496-4807
 * Department of Computer Science, Montana Tech
 */
const express = require('express'), // express module
    router = express.Router(), // define a route class
    fs = require('fs'), // filesystem routines
    bodyParser = require('body-parser'), // parse and post body routines
    strategy = require('./bsStrategy'), // tile selecting strategy
    http = require('http'); // for making http client requests

// log a structured message to the console for debugging purposes
function debugMsg(method, resource, ...param) {
    console.log(`[${method}]${resource} : ${param}`);
}

// Export the return value of this anonymous function
// pass in a function as a parameter that allows methods
// in this route class to send a message via SSE to its client
module.exports = function(sseMessage) {
    // parse the application/x-www-form-urlencoded
    router.use(bodyParser.urlencoded({ extended: true }));

    // parse the application/json
    router.use(bodyParser.json());

    // if filename parameter, associate with the req object
    router.param('filename', (req, res, next, filename) => {
        req.bsStateFn = filename;
        // continue processing middleware
        return next();
    }); // router.param

    // if url parameter, associate with the req object
    router.param('url', (req, res, next, url) => {
        req.opponentURL = url;
        // continue processing middleware
        return next();
    }); // router.param

    // if latency parameter, associated with the req object
    router.param('latency', (req, res, next, latency) => {
        req.latency = latency;
        // continue processing middleware
        return next();
    }); // router.param

    // if sessionID parameter, associate with the req object
    router.param('sessionID', (req, res, next, sessionID) => {
        req.sessionID = sessionID;
        // continue processing middleware
        return next();
    }); // router.param

    // Enter the battle phase of the game
    router.get('/battle/:filename/:url?', (req, res) => {
        // attempt to load the game model from bsStateFn
        // if not found return a 404 to the client via res object

        // load the bsState game model as the current model
        // change a global variable (or equiv) indicating the game is 
        // in battle mode - so you can reject additional requests

        // if opponentURL is present, use the http object to request a session
        // of the opponent - properly process the possible response
        // messages from your opponent request

        debugMsg('GET', '/battle', `${req.bsStateFn}`, `${req.opponentURL}`);
        res.sendStatus(200);
    }); // end /battle end-point

    // end point for accepting requests to enter a battle session with an opponent
    router.get('/session/:latency?', (req, res) => {
        // return to the requester a 412 response via the res object
        // if not in battle mode

        // return to the requester a 403 response via the res object
        // if already in a battle session with a different opponent

        // return to the requester a 400 response via the res object
        // if latency is provided and not within valid range

        // return to the requester a 200 response via the res object
        // generate a session identifier
        // compute the roll
        // set names for system and player
        // set the epoc
        // set the latency 

        // if we roll first, consult our strategy to generate a candidate tile
        // formulate a target message with the tile and the sessionID
        // invoke a /target request to the opponent using the http object and
        // change the roll to indicate opponent is next move

        console.log(strategy.getTile());
        debugMsg('GET', '/session', `${req.latency}`);
        res.sendStatus(200);
    }); // end /session end-point

    // accept a request from your opponent to prematurely terminate a game session
    router.delete('/session/:sessionID', (req, res) => {
        // return to the requester a 412 response via the res object
        // if not in battle mode

        // return to the requester a 403 response via the res object
        // if already in a battle session with a different opponent

        // return a 404 to the requester via the res object if the 
        // provided sessionID is incorrect 

        // return a 200 response to the requester via the res object 
        // and terminate the current session and reset all game state
        // note - will need to send a message via sseMessage to reset the 
        //        client - this you get to decide

        debugMsg('DELETE', '/session', `${req.sessionID}`);
        res.sendStatus(200);
    }); // end /session end-point

    // accept a target request from your opponent
    router.post('/target', (req, res) => {
        // return to the requester a 412 response via the res object
        // if not in battle mode

        // return to the requester a 403 response via the res object
        // if already in a battle session with a different opponent

        // return to the requester a 403 response via the res object
        // if it is not the opponent's turn

        // consult and update the game model based on the tile provided
        // send a message to the client via sseMessage
        // return a response sent to the opponent via the res object that corresponds 
        // to the outcome of consulting the game state model and compliant with the protocol

        // wait a latency amount of time and flip the roll of who's turn it is
        // using our strategy, generate a candidate tile, formulate a target req message
        // using the http object, send the target response
        // process the response from the opponent
        // update game state model to reflect the response from the opponent
        // update the client via sseMessage to refect the game model changes
        // check to see if we win and send a message indicating such to client via sseMessage
        // - if we win, reset the game state
        // change the roll to be the opponent's turn

        debugMsg('POST', '/target', req.body.sessionID, req.body.tile);
        res.sendStatus(200);
    }); // end /target end-point

    // return the route defined route class
    return router;
};