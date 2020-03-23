/**
 * S19 CSCI 470 Web Science
 * Authenticate user against private passphrase
 * using express-session
 * 
 * Phillip J. Curtiss, Assistant Professor
 * pcurtiss@mtech.edu, 406-496-4807
 * Department of Computer Science, Montana Tech
 */

const fs = require('fs'), // filesystem routines
    url = require('url'), // parsing and url routines
    bodyParser = require('body-parser'), // parse and post body routines
    express = require('express'), // express routines
    router = express.Router(); // custom route class 

// parse the application/x-www-form-urlencoded
router.use(bodyParser.urlencoded({ extended: true }));

// parse the application/json
router.use(bodyParser.json());

// If there is a filename parameter, attach
// filename to request object
router.param('filename', (req, res, next, filename) => {
    req.bsStateFn = filename;

    return next();
}); // end param

// end-point /bsState/:filename
// for methods - get, post, patch, delete
router.route('/:filename')
    .all((req, res, next) => {
        let path = url.parse(req.originalUrl).path;

        // redirect if bsStates encountered
        if (path.slice(0, 9) == '/bsStates') {
            res.redirect('/bsState' + path.slice(9, path.length));
            return;
        }

        // check for session id and validate, redirect to challenge
        // user to authenticate if needed - set bounce to our end-point
        if (!req.session || !req.session.clientKey) {
            res.redirect('/auth?bounce=' + path);
            return;
        }

        return next();
    })
    .get((req, res) => {
        console.log('Retrieve json from bsState: ' + req.bsStateFn);
        res.sendStatus(200);
    })
    .post((req, res) => {
        console.log('Save json to new bsState: ' + req.bsStateFn);
        res.sendStatus(200);
    })
    .patch((req, res) => {
        console.log('Save json to existing bsState: ' + req.bsStateFn);
        res.sendStatus(200);
    })
    .delete((req, res) => {
        console.log('Delete bsState: ' + req.bsStateFn);
        res.sendStatus(200);
    });

// end-point /bsStates
// for methods - get
router.route('/')
    .all((req, res, next) => {
        let path = url.parse(req.originalUrl).path;

        // redirect if bsState encountered
        if (path == '/bsState') {
            console.log('redirecting to /bsStates');
            res.redirect('/bsStates');
            return;
        }

        // check for session id an validate, redirect to challenge
        // user to authenticate if needed - set bounce to our end-point
        if (!req.session || !req.session.clientKey) {
            res.redirect('/auth?bounce=' + path);
            return;
        }

        return next();
    })
    .get((req, res) => {
        console.log('Would retrieve the files from the bsState directory.');
        res.sendStatus(200);
    });

module.exports = router;