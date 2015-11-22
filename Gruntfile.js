module.exports = function(grunt){

    grunt.initConfig({

        pkg: grunt.file.readJSON('package.json'),

        requirejs: {
            compileMain:{
                options: {
                    mainConfigFile: 'build/client/conf.js',
                    include: '../main',
                    findNestedDependencies: true,
                    optimize: 'none',
                    out: 'build/client/main.js',
                    onModuleBundleComplete: function (data) {
                        var fs = require('fs'),
                            amdclean = require('amdclean'),
                            outputFile = data.path;

                        fs.writeFileSync(outputFile, amdclean.clean({
                            'filePath': outputFile
                        }));
                    }
                }
            }
        },

        htmlmin: {
            build: {
                options: {
                    removeComments: true,
                    collapseWhitespace: true
                },
                files: [{
                    expand: true,
                    cwd: 'build/client/component',
                    src: '**/*.html',
                    dest: 'build/client/component'
                }]
            }
        },

        'json-minify': {
            build: {
                files: 'build/client/**/*.json'
            }
        },

        exec: {
            buildServer: {
                cmd: 'go build -o src/server/server.exe -v src/server/server.go'
            },
            startDevServer: {
                cmd: 'cd src/server && ./server.exe'
            },
            startBuildServer: {
                cmd: 'cd build/server && ./server.exe'
            },
            compileDevSass: {
                cmd: 'node node_modules/node-sass/bin/node-sass --output-style compressed -r -o src/client src/client'
            },
            compileBuildSass: {
                cmd: 'node node_modules/node-sass/bin/node-sass --output-style compressed -r -o build/client build/client'
            },
            watchDevSass: {
                cmd: 'node node_modules/node-sass/bin/node-sass --output-style compressed -w -r -o src/client src/client'
            },
            updateSeleniumServer: {
                cmd: 'node node_modules/protractor/bin/webdriver-manager update'
            },
            startSeleniumServer: {
                cmd: 'node node_modules/protractor/bin/webdriver-manager start'
            },
            testClient: {
                cmd: 'cd test/unit && node ../../node_modules/karma/bin/karma start'
            },
            testClientCI: {
                cmd: 'cd test/unit && node ../../node_modules/karma/bin/karma start --browsers PhantomJS'
            },
            testE2e: {
                cmd: 'cd test/e2e && node ../../node_modules/protractor/bin/protractor protractor.conf.js'
            }
        },

        copy: {
            serverExe: {
                src: 'src/server/server.exe',
                dest: 'build/server/server.exe',
                options : {
                    mode: true
                }
            },
            confJson: {
                src: 'src/server/conf.json',
                dest: 'build/server/conf.json',
                options : {
                    mode: true
                }
            },
            fullClient: {
                cwd: 'src/client',
                src: '**',
                dest: 'build/client/',
                expand: true,
                options : {
                    mode: true
                }
            }
        },
		
		jshint: {
			client: ['src/client/**/*.js',
						'!src/client/lib/**/*.js',
                        '!src/client/resource/**/*.js'],
			options: {
				asi: true
			}
		},

        processhtml: {
            clientIndex: {
                files: {
                    'build/client/index.html': ['build/client/index.html']
                }
            }
        },

        uglify: {
            mainJsBuild: {
                files: {
                    'build/client/main.js': ['build/client/main.js']
                }
            }
        },

        clean: {
            allClientBuildExcept_index_robot_favicon_main_legacy_resources: ['build/client/**/*', '!build/client/index.html', '!build/client/main.js', '!build/client/robots.txt', '!build/client/favicon.ico', '!build/client/resource/**', '!build/client/legacy/**'],
            buildCss: ['build/client/**/*.css', '!build/client/resource/**/*.css'],
            server: ['build/server', 'src/server/server.exe'],
            clientBuild: ['build/client'],
            clientTest: ['test/unit/coverage/*','test/unit/results/*'],
            sass: ['src/client/**/*.css'],
            e2e: ['test/e2e/results/*']
        },

        compass: {
            dev: {
                options: {
                    outputStyle: 'compressed',
                    cssPath: 'src/client',
                    sassPath: 'src/client',
                    watch: true,
                    trace: true
                }
            },
            build: {
                options: {
                    outputStyle: 'compressed',
                    cssPath: 'build/client',
                    sassPath: 'build/client',
                    watch: false,
                    trace: true
                }
            }
        }

    });

    grunt.loadNpmTasks('grunt-contrib-requirejs');
    grunt.loadNpmTasks('grunt-exec');
    grunt.loadNpmTasks('grunt-contrib-copy');
    grunt.loadNpmTasks('grunt-processhtml');
    grunt.loadNpmTasks('grunt-contrib-uglify');
    grunt.loadNpmTasks('grunt-contrib-clean');
    grunt.loadNpmTasks('grunt-contrib-compass');
    grunt.loadNpmTasks('grunt-contrib-htmlmin');
    grunt.loadNpmTasks('grunt-json-minify');
    grunt.loadNpmTasks('grunt-contrib-jshint');


    grunt.registerTask('buildServer', ['exec:buildServer', 'copy:serverExe', 'copy:confJson']);
    grunt.registerTask('cleanServer', ['clean:server']);

    grunt.registerTask('buildClient', ['jshint:client', 'copy:fullClient', 'clean:buildCss', 'exec:compileBuildSass', 'htmlmin:build', 'json-minify:build', 'requirejs:compileMain', 'uglify:mainJsBuild', 'processhtml:clientIndex', 'clean:allClientBuildExcept_index_robot_favicon_main_resources']);
    grunt.registerTask('testClient', ['exec:testClient']);
    grunt.registerTask('testClientCI', ['exec:testClientCI']);
    grunt.registerTask('cleanClientBuild', ['clean:clientBuild']);
    grunt.registerTask('cleanClientTest', ['clean:clientTest']);

    grunt.registerTask('buildAll', ['buildServer', 'buildClient']);
    //grunt.registerTask('buildAllWithCompass', ['buildServer', 'buildClientWithCompass']);
    grunt.registerTask('cleanAllBuild', ['cleanServer', 'cleanClientBuild']);

    grunt.registerTask('watchSass', ['exec:compileDevSass', 'exec:watchDevSass']);
    //grunt.registerTask('watchSassWithCompass', ['compass:dev']);
    grunt.registerTask('cleanSass', ['clean:sass']);

    grunt.registerTask('startDevServer', ['exec:startDevServer']);
    grunt.registerTask('startBuildServer', ['exec:startBuildServer']);

    grunt.registerTask('updateSeleniumServer', ['exec:updateSeleniumServer']);
    grunt.registerTask('startSeleniumServer', ['exec:startSeleniumServer']);

    grunt.registerTask('testE2e', ['exec:testE2e']);
    grunt.registerTask('cleanE2e', ['clean:e2e']);

    grunt.registerTask('nuke', ['cleanAllBuild', 'cleanClientTest', 'cleanSass', 'cleanE2e']);
	
    grunt.registerTask('lint', ['jshint:client']);

};
