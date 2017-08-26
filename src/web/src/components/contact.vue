<template>
    <div>
        <div class="contact">
            <h3>Contact</h3>
            <br>
            <p>Questions, issues or suggestions?</p>
            <br>
            <dt>Email</dt>
            <dd>cixx@cyberdyne.sh</dd>
            <dt>IRC</dt>
            <dd>trippysalmon @ freenode</dd>
            <dt>Slack</dt>
            <dd>trippysalmon @ bitcoincore.slack.com</dd>
            <dt>Twitter</dt>
            <dd>
                <a href="https://twitter.com/trippysalmon">@trippysalmon</a>
            </dd>
            <br>
            <br>
            <h3>Donations</h3>
            <br>
        </div>
        <div class="row justify-content-around text-center">
            <div class="col-md-6 col-sm-12">

                <p>Donate Bitcoin:</p>
                <p>1F43wkTWaUVyDbtfJrHXpc3C1ek5kcK8A</p>
                <a href="bitcoin:1F43wkTWaUVyDbtfJrHXpc3C1ek5kcK8A">
                    <img src="/static/qr_btc.png" alt="">
                </a>
                <br>
                <small>(image is clickable)</small>
            </div>

            <div class="col-md-6 col-sm-12">
                <div class="hidden-md-up">
                    <br>
                    <br>
                </div>
                <p>Donate Bitcoin Cash:</p>
                <p>1BFeA4phuDxfh4nqnRk9tjzYt9viQEvQkq</p>
                <br>
                <br>
                <img src="/static/qr_bch.png" alt="">
                <br>
                <br>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    name: 'contact'
}
</script>
<style scoped>
.contact {
    padding: 24px;
}
</style>

