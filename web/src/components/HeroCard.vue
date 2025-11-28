<template>
  <div class="hero min-h-dvh">
    <div class="hero-content m-0 p-0">
      <div class="bg-base-100 p-2 md:p-5 w-screen max-w-xl h-dvh max-h-200 flex flex-col justify-between gap-10">
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
