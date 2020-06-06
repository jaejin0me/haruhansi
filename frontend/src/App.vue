<template>
  <div id="app">
	  <Header></Header>
	  <Navbar></Navbar>
	  <Footer></Footer>
	  <!--
	  <Input v-on:addTodo="addTodo"></Input>
	  <List v-bind:propsdata="todoItems" @removeTodo="removeTodo"></List>
	  <Footer v-on:removeAll="clearAll"></Footer>
	  !-->
  </div>
</template>

<script>
import Footer from './components/Footer.vue'
import Navbar from './components/Navbar.vue'
import Input from './components/Input.vue'
import List from './components/List.vue'
import Header from './components/Header.vue'
import Mainpage from './components/Mainpage.vue'

export default {
	data() {
		return {
			todoItems: []
		}
	},
	methods: {
		addTodo(todoItem) {
			localStorage.setItem(todoItem, todoItem);
			this.todoItems.push(todoItem);
		},
		clearAll() {
			localStorage.clear();
			this.todoItems = [];
		},
		removeTodo(todoItem, index) {
			localStorage.removeItem(todoItem);
			this.todoItems.splice(index, 1);
		},
	},
	created() {
		if (localStorage.length > 0) {
			for (var i = 0; i < localStorage.length; i++) {
				this.todoItems.push(localStorage.key(i));
			}
		}
	},
	components: {
		'Header': Header,
		'Footer': Footer,
		'List': List,
		'Input': Input,
		'Navbar': Navbar,
	}
}
</script>

<style>
	body {
		text-align: center;
		background-color: #f6f6f8;
	}
	input {
		border-style: groove;
		width: 200px
	}
	button {
		border-style: groove;
	}
	.shadow {
		box-shadow: 5px 10px 10px rgba(0, 0, 0, 0.03);
	}
</style>
