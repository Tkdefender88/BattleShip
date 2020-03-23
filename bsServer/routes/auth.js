/**
 * S19 CSCI 470 Web Science
 * Authenticate user against private passphrase
 * using express-session
 * 
 * Phillip J. Curtiss, Assistant Professor
 * pcurtiss@mtech.edu, 406-496-4807
 * Department of Computer Science, Montana Tech
 */

const bodyParser = require('body-parser'), // parse http body part
    express = require('express'), // express module
    router = express.Router(), // define a route class
    crypto = require('crypto'), // used for hmac sh256 one-way hash
    privKey = 'Orange is the new crazy', // sha512 key

    // encrypted secret phrase for access to secure end-points
    phrase = '7fec7a2b065f38810b7df212dfd84bc90f08138309f93fe3935863a516eca220d4569371b245500286f8ef8905f8b0892b9b9deae2e2aaa64cade3a3e5e01ff8'; // our passphrase

// parse the application/x-www-form-urlencoded
router.use(bodyParser.urlencoded({ extended: true }));

// parse the application/json
router.use(bodyParser.json());

// default end-point for authentication
router.route('/')
    // if get method - render the auth form and remember the link that sent us here
    .get((req, res) => {
        res.render('auth', {
            title: 'Battleship Game Model Authentication',
            bounce: req.query['bounce'] || '/'
        }); // end render auth
    }) // end get method
    // if rest of form submission - process form fields
    .post((req, res) => {
        // save the referer from the post or /bsStates
        let referer = req.body.referer || '/bsStates';

        // create a new hash for update(write)/diget(read)
        let hash = crypto.createHash('sha512', privKey);

        // if passphrase agrees, then set session clientKey and 
        // redirect to our referer
        if (req.body.passphrase && req.body.passphrase &&
            hash.update(req.body.passphrase).digest('hex') == phrase) {
            req.session.clientKey = true;
            res.redirect(referer);
            return;
        } // end if

        // if the passphrase does not agree, render the auth form again
        // set the hidden bounce input elemnt to the referer
        res.render('auth', {
            title: 'Battleship Game Model Authentication',
            bounce: referer
        }); // end render auth
    }); // end post method

module.exports = router;