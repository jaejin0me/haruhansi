<template>
	<div class="Mainpage">
		{{ title }}
		{{ author }}
		{{ content }}
		<button v-on:click="fetchData">get data</button>
	</div>
</template>

<script>
import axios from 'axios';
//axios.defaults.headers.common['Access-Control-Allow-Origin'] = '*';

export default {
	data: function() {
		return {
			title: 'title',
			author: 'author',
			content : "content",
		}
	},
	methods: {
		fetchData: function() {
			axios.get('http://127.0.0.1:3000/poems/오르텅스 블루/사막')
			.then((response) => {
				console.log(response);
				console.log(response.data[0].content);
				this.content = response.data[0].content;
			})
			.catch((error) => {
				console.log(error);
			});
	    }
	}
}
</script>

<style>
	.Mainpage {
		color: #2f3b52;
		font-weight: 900;
		margin: 2.5rem 0 1.5rem
	}
</style>