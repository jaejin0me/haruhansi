<template>
	<div class="Mainpage">
		<p v-if=false>{{ id }}</p>
		<p>{{ title }}</p>
		<p>{{ author }}</p>
		<vue-markdown :source="content"></vue-markdown>
	</div>
</template>

<script>
import axios from 'axios';

export default {
	data: function() {
		return {
			id: 'empty',
			title: 'title',
			author: 'author',
			content : "content",
		}
	},
	created() {
		this.fetchData()
	},
	methods: {
		fetchData: function() {
			axios.get('http://haruhansi.com:3000/apoem/' + this.id)
			.then((response) => {
				this.id = response.data.id;
				this.content = response.data.content;
				this.title = response.data.title;
				this.author= response.data.author;
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