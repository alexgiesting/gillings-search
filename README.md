# gillings-search
A web app for querying research publications at UNC Gillings. Web app can be found at http://gillings-search-giesting.apps.cloudapps.unc.edu/


## Gillings Search Tool

Our team is developing a web-based tool for querying research publications at the[  Gillings School of Global Public Health](https://sph.unc.edu/). This tool will help showcase the contributions of researchers at our school to legislators, journalists, students, and scientists from around the world.

##  Getting started

Before developing or running this software, there are a few softwares which need to be installed onto your machine. These include the Go, Node.js, and docker extension to vscode. Any developer working on the project can install the Go extension from the vs code[  market place](https://marketplace.visualstudio.com/items?itemName=golang.go) and should work once the extension is installed. To utilize React, Node.js must first be installed on the user's computer and can be found at this link[  here](https://nodejs.org/en/download/)  . It is important that you download the  correct version for the operating system that you will be using. Once this is done, React should be fully functional and a new React app can be created using the command "`npx create-react-app my-app`". Finally, to utilize Docker it first must be[  installed onto your machine](https://docs.docker.com/get-docker/) and added to the system path. On linux, you should[  also enable Docker CLI for the non-root user account](https://docs.docker.com/engine/install/linux-postinstall/) that will be used to run Vscode. Then, the user can install the[  docker extension for Vscode](https://marketplace.visualstudio.com/items?itemName=ms-azuretools.vscode-docker) to provide future support and from here it should be functional. To set up the development environment with docker, a user will simply need to navigate to the root directory, run "`docker compose build`" and then run "`docker compose up`". This activates the docker containers to allow for a working development environment to run. Once everything is installed and the containers are running, a user can start a local server using npm, running "`npm start`". This command starts up a local server on which the project will be hosted. These instructions were tested and verified to work by Andrew Chun and Alex Lewis on Windows computers, but the same process will work for other operating systems.

##   Testing

Currently the test suite for our source code is broken up into two parts: the backend portion which is tested using the built in testing functionality of the Go languages and the front end portion which is tested using the React testing library. In order to run the test suite for the backend, the user must create test files whose names end with "`_test.go`" and  then run the command "`go test`". As of now, the code should also automatically run the test.go files as well when the `docker-compose` command is used. This will run the test codes in the current directory  and provide test coverage percentages, time taken, and error messages if needed. To test the front end with the React testing library, it is done via jest. The naming conventions for these files will end in the suffix `.test.js` and can be run by using the command npm test in the terminal while in the react directory and can see test coverage with `npm test -- --coverage`. As of now there is no different command for unit tests vs integration tests.

##  Deployment

This project's codebase is hosted on a GitHub repository, to which a new developer would need to gain access to by talking to one of the current team members. Currently, there are no dedicated staging or pre-production environments at the time of this documentation's writing. However, this will be developed in the future and utilized. The fully deployed software in the future would use a backend managed by Carolina cloud apps utilizing its mongoDB and openshift capabilities while pulling citations from the Gillings' Scopus database and querying them with Solr with the frontend mainly utilizing react.js.Continuous integration is enabled through use of a central GitHub repo and following git practices. Continuous deployment is also enabled through docker containers but there aren't any automated tests at this moment. 

##  Technologies used

The technologies chosen for this project include some which have been provided by the university services in order to be easily integrated to the current Gillings system which remains connected to UNC Chapel Hill services. This includes Carolina cloud apps and OpenShift for hosting services, Scopus for providing citations as provided by Gillings, and WordPress to host our team's personal website. Other technologies incorporated into our tool include GO to construct our backend, Mongo DB in order to host our database, React.js for our frontend network, Solr for load balanced querying, Docker for server hosting, and excel and Zotero for external APIs and citation management. Links to the technologies can be found below:

1.  Carolina Cloud Apps: <https://cloudapps.unc.edu/> 

2.  OpenShift: <https://www.redhat.com/en/technologies/cloud-computing/openshift> 

3.  Scopus: <https://www.scopus.com/search/form.uri?display=basic> 

4.  WordPress: <https://wordpress.com/> 

5.  Go: <https://golang.org/> 

6.  MongoDB: <https://www.mongodb.com/> 

7.  React.js: <https://reactjs.org/>

8.  Solr:[  https://solr.apache.org/](https://solr.apache.org/)

9.  Docker: <https://www.docker.com/> 

10. Excel: <https://www.microsoft.com/en-us/microsoft-365/excel> 

11. Zotero:[  https://www.zotero.org/](https://www.zotero.org/)

A copy of our architecture decision record is available at this link:[  https://tarheels.live/comp523f21o/deliverables/application-architecture/](https://tarheels.live/comp523f21o/deliverables/application-architecture/)

As well as in the ADR.md, found in the top level of the main branch.

##  End of Semester Overview

Here is a video overview of all that we accomplished by the end of the 2021 fall academic semester, which was the original timeframe that we intended to work on the project. The original creators may continue to work on the project after this timeframe. 

https://user-images.githubusercontent.com/49577025/147509733-385c5b9f-3ffd-4939-b783-676c983aa90f.mp4

##  Contributing

In order for a new developer to get started on the project and make contributions, there are a few accesses they need to obtain. Firstly, they will need access to the GitHub repo and the Trello so they may have access to the code and the project schedule. Next, they need access to Carolina cloud apps in order to make any changes to the hosting services if needed. Otherwise, the README that our team has developed should allow for a straightforward onboarding process. In terms of styping, our team uses the Prettier formatter (<https://marketplace.visualstudio.com/items?itemName=esbenp.prettier-vscode>) which is an extension from the VS code marketplace to stylize our code. If you would like more information about the project this can be found at the project website: <https://tarheels.live/comp523f21o/>

##  Authors

The major authors on this project are:

Andrew Chun - achun@live.unc.edu -- Client Manager, Project Manager

Alex Lewis - giesting@live.unc.edu -- Webmaster

Justin Bautista -- Justinbb@live.unc.edu -- Tech Lead

Together we have established and continued to work on developing the web based querying tool for the Gillings school of Global Public Health and have kept our source code updated to reflect the progress that we have made.

## License

All of the source code in this repository and other creative works that may be related to this project are property of the University of North Carolina at Chapel Hill as well as the Gillings School of Public Health. These are protected under the GNU General Public License v3.0 so that way our client has the legal right to commercial use, modification, distribution, private use, and other methods of which they wish to use and extend our work in whatever way they would like. 

## Acknowledgements

Team O would like to formally acknowledge the valuable guidance and support of our mentor, Ben Knoble along side our COMP 523 Professor Dr. Jeff Terrell for providing important insight and knowledge about the process of software development. We would also like to extend gratitude towards Alexia Kelly, Penny Gordon-Larsen, Adam Dodd, Paul Glass-Steel, and anyone else at the Gillings School of Global Public Health who may have helped us along this project for their advice and resources regarding their current database management system.
