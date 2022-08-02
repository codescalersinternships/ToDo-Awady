<script>
    import TodoItem from './TodoItem.svelte';
    import { onMount } from 'svelte'

    let text = '';
    let currentFilter = 'all';
    let ServerURL = 'http://127.0.0.1:8080/todo'
    let todos = [];
    onMount(async() => {
		const res = await fetch(ServerURL, {
			method: 'GET',
		})
		todos = await res.json()
	})
    async function getToDo() {
        const res = await fetch(ServerURL, {
			method: 'GET',
		})
		todos = await res.json()
    }
    async function addTodo(event) {
        if (event.key === 'Enter') {
            todos = [...todos, {
                id: (todos.length === 0) ? 1 : todos[todos.length - 1].id + 1,
                completed: false,
                text: text
            }];
            const res = await fetch(ServerURL, {
                method: 'POST',
                mode: 'no-cors',
                body: JSON.stringify({
                    text,
                })
		    })
        }
    }

    $: todosRemaining = filteredTodos.filter(todo => !todo.status).length;
    $: filteredTodos = currentFilter === 'all' ? todos : currentFilter === 'completed'
        ? todos.filter(todo => todo.status)
        : todos.filter(todo => !todo.status)



    async function handleDeleteTodo(event) {
		await fetch(ServerURL + '/' + event.detail.id, {
			method: 'DELETE',
		}).then(() => {
            getToDo()
        })
    }

    async function handleToggleComplete(event) {
        await fetch(ServerURL + '/' + event.detail.id, {
                method: 'PUT',
                body: JSON.stringify({
                    "text" : event.detail.text,
                    "status" : event.detail.status
                }),
            }).then(() => {
                getToDo()
        })
    }
</script>

<style>
    .container {
        max-width: 800px;
        margin: 10px auto;
    }
    .todo-input {
        width: 100%;
        padding: 10px, 20px;
        font-size: 18px;
        margin-bottom: 20px;
    }
    .inner-container {
        display: flex;
        align-items: center;
        justify-content: space-between;
        font-size: 16px;
        border-top: 1px solid lightgrey;
        padding-top: 15px;
        margin-bottom: 13px;
    }
</style>

<div class="container">

    <h2>Svelte Todo App</h2>
    <input type="text" class="todo-input" placeholder="Insert todo item ..." bind:value={text} on:keydown={addTodo} >

    {#each filteredTodos as todo}
        <div class="todo-item">
            <TodoItem {...todo} on:deleteTodo={handleDeleteTodo} on:toggleComplete={handleToggleComplete} />
        </div>
    {/each}

    <div class="inner-container">
        <div>{todosRemaining} items left</div>
    </div>
</div>