// Template -- To replace

<script setup lang="ts">
import { ref, reactive } from 'vue'

const showPassword = ref(false)

const form = reactive({
  email: '',
  password: '',
  rememberMe: false,
})

function handleLogin() {
  console.log('form:', form)
  // TODO: call API
}
</script>

<template>
  <div class="login-page">
    <div class="login-card">
      <h1>Connexion</h1>

      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label for="email">Email</label>
          <input id="email" v-model="form.email" type="email" placeholder="vous@exemple.com" />
        </div>

        <div class="form-group">
          <div class="label-row">
            <label for="password">Mot de passe</label>
            <router-link to="forgot-password">Mot de passe oublié ?</router-link>
          </div>
          <div class="password-wrapper">
            <input
              id="password"
              v-model="form.password"
              :type="showPassword ? 'text' : 'password'"
              placeholder="••••••••"
            />
            <button type="button" @click="showPassword = !showPassword">
              {{ showPassword ? 'Masquer' : 'Afficher' }}
            </button>
          </div>
        </div>

        <div class="form-group">
          <label>
            <input type="checkbox" v-model="form.rememberMe" />
            Se souvenir de moi
          </label>
        </div>

        <button type="submit">Se connecter</button>
      </form>
    </div>
  </div>
</template>
