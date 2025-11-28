<template>
  <div class="hero bg-base-200 min-h-dvh">
    <div class="hero-content h-full m-0 p-0">
      <div class="bg-base-100 p-2 md:p-5 w-screen max-w-xl flex flex-col justify-between h-full max-h-dvh sm:max-h-200 gap-10">
        <header>
          <slot name="header" />
        </header>
        <main class="flex-1">
          <slot />
        </main>
        <footer>
          <slot name="footer" />
        </footer>
      </div>
    </div>
  </div>
</template>
