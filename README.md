1. Purpose:
   
   This Backend API purposed to fulfill amartha backend test held in September 2024
   i was choosing task NUMBER ONE since it was core part of loan system. But if you need later, i can work on another task number as well.
2. Documentation
   
   - Api Documentation was served by postman server in [https://documenter.getpostman.com/view/3715067/2sAXjF6tr4](https://documenter.getpostman.com/view/3715067/2sAXqnePxJ)
     or can be downloaded here https://drive.google.com/file/d/17qvqaNmbDmyVfbD5npHnJeTrYlW3LXvI/view?usp=sharing
   - ERD Documentation :
     <img width="926" alt="Screenshot 2024-09-13 at 08 54 46" src="https://github.com/user-attachments/assets/43919570-52d6-4afc-a673-9dc2d6f98e9a">

3. How to run
   - install go.
   - git clone https://github.com/DeniesKresna/amarthatest.git
   - run docker compose up -d for starting the database
   - run go run main.go to start develop and testing
     
  if you re using docker compose and only want to test, i have setup the dockerfile for the app, adjust it to match your pc, since im not using it and comment the docker compose on app section when i developed this task

4. Technology Background
   
   I was using golang as required by the test, and using GIN as web framework since GIN has gain the most stars as go web framework in github and based on its performance as the fastest go webframework.
   I was using MySql v8 for the database
   I did not use clean architecture since this only for doing small part of big service. this implement simple structure, (controller-repo)
   I was using cron as well

5. Disclaimer
   - I did not using env for this task
   - I cannot test the 7 days cron, so i had put comment to help you understand my logic flow
   - I didnt add any other function to manage users, authentication, etc.
   - This code only small part of big service, can be improve a lot for making good service. i only share my logic think for the test task

   
