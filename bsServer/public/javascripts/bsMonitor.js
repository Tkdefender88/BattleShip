/**
 * S19 CSCI470 Web Science
 * 
 * SSE Behaviors for the SSE Events
 *
 * Phillip J. Curtiss, Assistant Professor
 * Department of Computer Science, Montana Tech
 * pcurtiss@mtech.edu, 406-496-4807
 * (c) 2018 All Rights Reserved
 */

(function(conf) {
    // HTML element is used as a proxy for adding event listeners and dispatching
    // those events which are custom for this application
    var htmlElement = document.querySelector('html'),

        // our custom event that will get dispatched whenever the game model is updated
        // this will cause the event listener(s) to trigger and update the information 
        // panels for each ship based on the new data in the model
        modelUpdate = new Event('modelUpdate'),

        // initial game model object
        gameModel = {},

        healthPercent = new Intl.NumberFormat('en-US', {
            style: 'percent',
            maximumFractionDigits: 2
        });

    // event source for SEEs
    evtSource = null;

    // create an information panel for each configured ship
    // by cloning the template and modifying attributes
    if (conf) {
        // template holds our information panel for each ship in the game
        let template = document.querySelector('[title="Ship Detail Information"]'),

            // obtain the main element node from the document
            mainElement = document.querySelector('main');

        // iterate over the configured ships
        conf.ships.forEach((ship) => {
            // create a clone (parentless) of the template
            let clone = template.content.cloneNode(true);

            // modify attributes of the cloned template
            // then add it to the main element node as a child
            if (mainElement && clone) {
                // obtain a reference to the article element node
                let articleElement = clone.childNodes[1];

                // modify the attributes and add clone to main element
                if (articleElement) {
                    // set the data-name attribute for this clone
                    articleElement.setAttribute('data-name', ship.name);

                    // set the data-size attribute for this clone
                    articleElement.setAttribute('data-size', ship.size);

                    // insert the clone into the main element node
                    mainElement.insertBefore(clone, mainElement.firstChild);
                }
            }
        });
    } else { // no document content to populate
        return;
    }

    // Event listener for modelUpdate events
    // update the information panels based on the current state of the gameModel
    htmlElement.addEventListener('modelUpdate', (event) => {
        // regex used to replace shipname with actual shipname from game Model
        let re = /(\[data-name=")shipname("\])/;

        // iterate over the ships in our configuration object and
        // update the corresponding text node of the output node elements on the page
        conf.ships.forEach((ship) => {

            if (!gameModel[ship.name]) {
                gameModel[ship.name] = {};
                gameModel[ship.name].name = ship.name;
                gameModel[ship.name].placement = ship.placement;
                gameModel[ship.name].hitProfiles = new Array(new Array(), new Array());
                gameModel[ship.name].shipSize = ship.size;
            }

            if (!gameModel.misses) {
                gameModel.misses = new Array();
            }

            // extract the data from the gameModel and store in local variables
            let reFmt = '$1' + ship.name + '$2';

            // iterate over the output selectors and populate the element node
            // with the data extracted from the game Model
            conf.selectors.forEach((selector) => {
                // obtain the element through applying the regex replace on the selector
                // from the conf object passed in as a parameter
                let element = document.querySelector(selector.value.replace(re, reFmt));

                // match the selector for which a value needs updating
                switch (selector.key) {
                    case 'shipName':
                        element.value = (gameModel[ship.name]) ? gameModel[ship.name].name : '';
                        element.value += ` (${gameModel[ship.name].shipSize})`;
                        element.value = element.value.charAt(0).toUpperCase() + element.value.slice(1);

                        if (gameModel[ship.name].shipSize == gameModel[ship.name].hitProfiles[0].length) {
                            element.value += ' <- SUNK ->';
                        }
                        break;
                    case 'gridTile':
                        element.value = gameModel[ship.name].placement.substr(0, 2);
                        break;
                    case 'orientation':
                        element.value = (gameModel[ship.name].placement.substr(2, 1) == 'H') ? 'Hoirz' : 'Vert';
                        break;
                    case 'playerTile':
                        element.value = gameModel[ship.name].hitProfiles[0].join(' ');
                        break;
                    case 'opponentTile':
                        element.value = gameModel[ship.name].hitProfiles[1].join(' ');
                        break;
                    case 'playerHealth':
                        element.value = healthPercent.format(1 - gameModel[ship.name].hitProfiles[0].length / gameModel[ship.name].shipSize);

                        if (gameModel[ship.name].shipSize == gameModel[ship.name].hitProfiles[0].length) {
                            element.parentNode.parentNode.classList.add('sunk');
                        }
                        break;
                    case 'opponentHealth':
                        element.value = healthPercent.format(1 - gameModel[ship.name].hitProfiles[1].length / gameModel[ship.name].shipSize);
                        break;
                }
            });
        });
        // update the text node of output element for misses with gameModel data
        document.querySelector(conf.misses).value = (gameModel.misses) ? gameModel.misses.join(' ') : '';
    });

    document.querySelector('#refresh').addEventListener('click', (event) => {
        gameModel = {};
        htmlElement.dispatchEvent(modelUpdate);


    });

    // Register listener for buttons
    document.querySelector('#connect').addEventListener('click', (event) => {
        if (evtSource != null) {
            evtSource.close();
            setTimeout(() => {
                evtSource = null;
                event.target.innerHTML = 'Connect';
            }, 3000);
        } else {
            evtSource = new EventSource(conf.gameSvr + '/stream');

            document.getElementById('refresh').disabled = true;

            // event listeners if there is an error encountered
            // with the stream 
            evtSource.addEventListener('error', () => {
                event.target.innerHTML = 'Connecting...';
            });

            // once the stream is open and data strarts flowing
            evtSource.addEventListener('open', () => {
                event.target.innerHTML = 'Disconnect';
            });

            // if there is some information sent that is not captured
            // by one of the custom events, then write to console
            evtSource.addEventListener('message', (evt) => {
                console.log(evt);
            });

            evtSource.addEventListener('END', (evt) => {
                let connectBtn = document.querySelector('#connect');

                evtSource.close();
                connectBtn.innerHTML = 'Connect';
            });

            // if the event is a miss event
            evtSource.addEventListener('MISS', (evt) => {
                gameModel.misses.push(evt.data);
                htmlElement.dispatchEvent(modelUpdate);
            });

            evtSource.addEventListener('DESTROYER', (evt) => {
                gameModel['destroyer'].hitProfiles[0].push(evt.data);
                htmlElement.dispatchEvent(modelUpdate);
            });

            evtSource.addEventListener('SUBMARINE', (evt) => {
                gameModel['submarine'].hitProfiles[0].push(evt.data);
                htmlElement.dispatchEvent(modelUpdate);
            });

            evtSource.addEventListener('CRUISER', (evt) => {
                gameModel['cruiser'].hitProfiles[0].push(evt.data);
                htmlElement.dispatchEvent(modelUpdate);
            });

            evtSource.addEventListener('BATTLESHIP', (evt) => {
                gameModel['battleship'].hitProfiles[0].push(evt.data);
                htmlElement.dispatchEvent(modelUpdate);
            });

            evtSource.addEventListener('CARRIER', (evt) => {
                gameModel['carrier'].hitProfiles[0].push(evt.data);
                htmlElement.dispatchEvent(modelUpdate);
            });
        }
    });


    // Cause an initial population (or clearing) of the output element nodes
    htmlElement.dispatchEvent(modelUpdate);
})({ // configuration object - links structure and behavior
    'ships': [ // Configured ships
        { 'name': 'carrier', 'size': '5', 'placement': 'J2H' },
        { 'name': 'battleship', 'size': '4', 'placement': 'H2H' },
        { 'name': 'cruiser', 'size': '3', 'placement': 'F3H' },
        { 'name': 'submarine', 'size': '3', 'placement': 'D4H' },
        { 'name': 'destroyer', 'size': '2', 'placement': 'B6H' }
    ],
    'gameSvr': 'https://csdept16.mtech.edu:30120/sse', // game server (backend)
    'selectors': [ // selectors for the output element nodes
        { 'key': 'shipName', 'value': '[data-name="shipname"]>div.shipName>output' },
        { 'key': 'gridTile', 'value': '[data-name="shipname"]>div.gridTile>output' },
        { 'key': 'orientation', 'value': '[data-name="shipname"]>div.orientation>output' },
        { 'key': 'playerTile', 'value': '[data-name="shipname"]>div.playerTile>output' },
        { 'key': 'opponentTile', 'value': '[data-name="shipname"]>div.opponentTile>output' },
        { 'key': 'playerHealth', 'value': '[data-name="shipname"]>div.playerHealth>output' },
        { 'key': 'opponentHealth', 'value': '[data-name="shipname"]>div.opponentHealth>output' }
    ],
    'misses': 'article>div#misses>output' // selector for the misses output element node
});