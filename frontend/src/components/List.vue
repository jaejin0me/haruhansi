<template>
	<section>
		<transition-group name="list" tag="ul">
			<li v-for="(todoItem, index) in propsdata" :key="todoItem" class="shadow">
				<i class="checkBtn fa fa-check" aria-hidden="true"></i>
				{{ todoItem }}
				<span class="removeBtn" type="button" @click="removeTodo(todoItem, index)">
					<i class="far fa-trash-alt" aria-hidden="true"></i>
				</span>
			</li>
		</transition-group>
	</section>
</template>

<script>
export default {
	props: ['propsdata'],
	methods: {
		removeTodo(todoItem, index) {
			this.$emit('removeTodo', todoItem, index)
		}
	}
}
</script>

<style>
	.list-item {
		display: inline-block;
		margin-right: 10px;
	}
	.list-move {
		transition: transform 1s;
	}
	.list-enter-active, .list-leave-active {
		transition: all 1s;
	}
	.list-enter, .list-leave-to {
		opacity: 0;
		transform: translateY(30px);
	}
</style>
