import 'bootstrap'
import $ from 'jquery';
import '@fortawesome/fontawesome-free/css/all.css'
import './index.scss' // Import our scss file

window.onload = () => {
    console.log('window loaded');
    $("#msgid").html("This is Hello World by JQuery");
};