<script setup>
import { ref, reactive } from 'vue'

new Vue({
  el: '#v-for-object',
  data: {
    object: {
      title: 'How to do lists in Vue',
      author: 'Jane Doe',
      publishedAt: '2016-04-10'
    }
  }
})

</script>

<template>
  <header>
    <div class="wrapper">
    </div>
  </header>


  <main>
    <div class="header">
      <h1>
        Budgeting App
      </h1>
    </div>
    <div  class="content">
      <div class="incomes">
        <div class="title">
          <h1>
            Incomes
          </h1>
        </div>
        <div class="listedItems">
          <ul id="v-for-object" class="demo">
            <li v-for="value in object">
              {{ value }}
            </li>
        </ul>
        </div>
        
  </main>
</template>

<style scoped>
header {
  line-height: 1.5;
}

.logo {
  display: block;
  margin: 0 auto 2rem;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }
}
</style>
