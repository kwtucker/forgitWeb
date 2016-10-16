var gulp = require('gulp'),
    sass = require('gulp-sass'),
    gulpgo = require("gulp-go"),
    concat = require('gulp-concat'),
    rename = require('gulp-rename'),
    uglify = require('gulp-uglify');

var go,

gulp.task('styles', function() {
    gulp.src('assets/scss/**/*.scss')
        .pipe(sass({outputStyle: 'compressed'}).on('error', sass.logError))
        // .pipe(gulp.dest('static/css/'))
        .pipe(rename('style.min.css'))
        .pipe(gulp.dest('static/css/'));
});

gulp.task('styles:watch', function () {
  gulp.watch('assets/scss/**/*.scss', ['styles']);
});

gulp.task('js:watch', function () {
  gulp.watch('assets/js/**/*.js', ['js']);
});

// Server side
gulp.task("go-run", function() {
  go = gulpgo.run("main.go", [], {cwd: __dirname, stdio: 'inherit'});
});

gulp.task("devs", ["go-run"], function() {
  gulp.watch([__dirname+"/**/*.go", __dirname+"/**/*.html"]).on("change", function() {
    go.restart();
  });
});

gulp.task('js', function() {
  return gulp.src('assets/js/**/*.js')
    .pipe(concat('main.js'))
    .pipe(rename('main.min.js'))
    .pipe(uglify())
    .pipe(gulp.dest('static/js/'));
});

gulp.task('startup', ['devs', 'styles', 'styles:watch', 'js', 'js:watch']);
