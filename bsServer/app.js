/**
 * S19 CSCI 470 Web Science
 * Authenticate user against private passphrase
 * using express-session
 * 
 * Phillip J. Curtiss, Assistant Professor
 * pcurtiss@mtech.edu, 406-496-4807
 * Department of Computer Science, Montana Tech
 */

const createError = require('http-errors'),
    express = require('express'),
    path = require('path'),
    cookieParser = require('cookie-parser'),
    session = require('express-session'),
    logger = require('morgan'),

    indexRouter = require('./routes/index'),
    usersRouter = require('./routes/users'),
    bsStateRouter = require('./routes/bsState'),
    bsAuth = require('./routes/auth'),
    sse = require('./routes/sse'),
    bsProtocol = require('./routes/bsProtocol')(sse.sseMessage);

app = express(),

    expTime = 60 * 60 * 1000,
    privKey = 'Orange is the new crazy';

// view engine setup
app.set('views', path.join(__dirname, 'views'));
app.set('view engine', 'jade');

// middleware used
app.use(logger('dev'));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(express.static(path.join(__dirname, 'public')));

// setup session 
app.use(session({
    secret: privKey,
    name: 'bsSession',
    resave: false,
    saveUninitialized: true,
    expires: Date.now() + expTime
}));

// establish route classes
app.use('/', indexRouter);
app.use('/users', usersRouter);
app.use('/bsStates?', bsStateRouter);
app.use('/auth', bsAuth);
app.use('/sse', sse.router);
app.use('/bsProtocol', bsProtocol);

// catch 404 and forward to error handler
app.use(function(req, res, next) {
    next(createError(404));
});

// error handler
app.use(function(err, req, res, next) {
    // set locals, only providing error in development
    res.locals.message = err.message;
    res.locals.error = req.app.get('env') === 'development' ? err : {};

    // render the error page
    res.status(err.status || 500);
    res.render('error');
});

module.exports = app;