/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

const gulp = require('gulp');
const del = require('del');
const replace = require('gulp-replace');
const gulpif = require('gulp-if');
const imagemin = require('gulp-imagemin');
const livereload = require('gulp-livereload');
const zip = require('gulp-zip');

const target = process.env.TARGET;

gulp.task('manifest', () => {
  const pkg = require('./package.json');

  return gulp
    .src(`manifests/${target}/manifest.json`)
    .pipe(replace('__VERSION__', pkg.version))
    .pipe(gulp.dest(`dist/${target}`));
});

gulp.task('styles', () => {
  return gulp.src('src/styles/*.css').pipe(gulp.dest(`dist/${target}/styles`));
});

gulp.task(
  'html',
  gulp.series('styles', () => {
    return gulp.src('src/*.html').pipe(gulp.dest(`dist/${target}`));
  })
);

gulp.task('images', () => {
  return gulp
    .src('src/images/**/*')
    .pipe(
      gulpif(
        gulpif.isFile,
        imagemin({
          progressive: true,
          interlaced: true,
          svgoPlugins: [{ cleanupIDs: false }]
        })
      )
    )
    .pipe(gulp.dest(`dist/${target}/images`));
});

gulp.task(
  'clean',
  del.bind(null, ['.tmp', `dist/${target}`, `package/${target}`])
);

gulp.task(
  'watch',
  gulp.series('manifest', 'html', 'styles', 'images', () => {
    livereload.listen();

    gulp
      .watch([
        'src/*.html',
        'src/scripts/**/*',
        'src/images/**/*',
        'src/styles/**/*'
      ])
      .on('change', livereload.reload);

    gulp.watch('src/*.html', gulp.parallel('html'));
    gulp.watch('manifests/**/*.json', gulp.parallel('manifest'));
  })
);

gulp.task('package', function() {
  const manifest = require(`./dist/${target}/manifest.json`);

  return gulp
    .src(`dist/${target}/**`)
    .pipe(zip('dnote-' + manifest.version + '.zip'))
    .pipe(gulp.dest(`package/${target}`));
});

gulp.task('build', gulp.series('manifest', gulp.parallel('html', 'images')));

gulp.task('default', gulp.series('clean', 'build'));
