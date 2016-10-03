var gulp = require('gulp'),
    sass = require('gulp-sass'),
    gulpgo = require("gulp-go");

var go;

gulp.task('styles', function() {
    gulp.src('scss/**/*.scss')
        .pipe(sass({outputStyle: 'compressed'}).on('error', sass.logError))
        .pipe(gulp.dest('static/css/'))
});

gulp.task('styles:watch', function () {
  gulp.watch('scss/**/*.scss', ['styles']);
});

// Server side
gulp.task("go-run", function() {
  go = gulpgo.run("main.go", [], {cwd: __dirname, stdio: 'inherit'});
});

gulp.task("devs", ["go-run"], function() {
  gulp.watch([__dirname+"/**/*.go", __dirname+"/**/*.html" ]).on("change", function() {
    go.restart();
  });
});

gulp.task('startup', ['devs', 'styles', 'styles:watch']);
