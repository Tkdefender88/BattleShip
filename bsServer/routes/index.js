var express = require('express');
var router = express.Router();

/* GET home page. */
router.get('/', function(req, res, next) {
  res.render('index', { title: 'Express' });
});

/* Get Save Game Model */
router.get('/push', function(req, res, next) {
   res.render('push', {title: 'Save Game Model'});
});

module.exports = router;
