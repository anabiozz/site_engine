var gulp = require('gulp');
var sass = require('gulp-sass');
var babel = require('gulp-babel');
var eslint = require('gulp-eslint');
var browserSync = require('browser-sync').create();

gulp.task('browserSync', function() {
	browserSync.init({
		server: {
			baseDir: 'app'
		},
	})
});

/*gulp.task('babel_node', function() {
	return gulp.src('app/js/node/es6/*.js')
	.pipe(babel())
	.pipe(eslint())
	.pipe(eslint.format())
	.pipe(gulp.dest('app/js/node/dist'))
});*/

gulp.task('babel', function() {
	return gulp.src('app/js/public/es6/*.js')
	.pipe(babel())
	.pipe(eslint())
	.pipe(eslint.format())
	.pipe(gulp.dest('app/js/public/dist'))
});

gulp.task('sass', function () {
	return gulp.src('app/sass/styles.scss')
	.pipe(sass())
	.pipe(gulp.dest('app/stylesheets'))
	.pipe(browserSync.reload({
		stream: true
	}))
});

/*gulp.task('autoprefixer', function () {
    var postcss      = require('gulp-postcss');
    var sourcemaps   = require('gulp-sourcemaps');
    var autoprefixer = require('autoprefixer');

    return gulp.src('app/stylesheets/*.css')
        .pipe(sourcemaps.init())
        .pipe(postcss([ autoprefixer({ browsers: ['last 2 versions'] }) ]))
        .pipe(sourcemaps.write('.'))
        .pipe(gulp.dest('app/stylesheets'));
});*/

gulp.task('watch', ['browserSync', 'sass', 'babel'], function() {
	gulp.watch('app/sass/_about.scss', ['sass']);
	gulp.watch('app/sass/_banner.scss', ['sass']);
	gulp.watch('app/sass/_base.scss', ['sass']);
	gulp.watch('app/sass/_contact.scss', ['sass']);
	gulp.watch('app/sass/_footer.scss', ['sass']);
	gulp.watch('app/sass/_header.scss', ['sass']);
	gulp.watch('app/sass/_projects.scss', ['sass']);
	gulp.watch('app/sass/_reset.scss', ['sass']);
	gulp.watch('app/sass/_variables.scss', ['sass']);
	gulp.watch('app/sass/styles.scss', ['sass']);
	gulp.watch('app/index.html', browserSync.reload);
	gulp.watch('app/js/public/es6/*.js', ['babel']);
});

gulp.task('default', function() {
	gulp.src(['app/js/public/es6/*.js', 'app/js/public/node/*.js'])
	.pipe(eslint())
	.pipe(eslint.format());
});