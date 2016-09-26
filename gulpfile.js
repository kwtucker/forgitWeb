var gulp = require('gulp'),
    sass = require('gulp-sass');


gulp.task('styles', function() {
    gulp.src('scss/**/*.scss')
        .pipe(sass({outputStyle: 'uncompressed'}).on('error', sass.logError))
        .pipe(gulp.dest('static/css/'))
});

gulp.task('styles:watch', function () {
  gulp.watch('scss/**/*.scss', ['styles']);
});

gulp.task('startup', ['styles', 'styles:watch']);
