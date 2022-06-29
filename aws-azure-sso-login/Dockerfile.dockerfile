FROM node:13.8-stretch

WORKDIR /tmp

RUN apt update \
    && apt install -y curl \
    && apt install -y gconf-service libasound2 libatk1.0-0 libc6 libcairo2 libcups2 libdbus-1-3 libexpat1 libfontconfig1 libgcc1 libgconf-2-4 libgdk-pixbuf2.0-0 libglib2.0-0 libgtk-3-0 libnspr4 libpango-1.0-0 libpangocairo-1.0-0 libstdc++6 libx11-6 libx11-xcb1 libxcb1 libxcomposite1 libxcursor1 libxdamage1 libxext6 libxfixes3 libxi6 libxrandr2 libxrender1 libxss1 libxtst6 ca-certificates fonts-liberation libappindicator1 libnss3 lsb-release xdg-utils wget \
    && apt install -y sudo \
    && rm -rf /var/lib/apt/lists/*

RUN npm install -g aws-azure-login --unsafe-perm

RUN \
    groupadd -g 999 awsuser && useradd -u 999 -g awsuser -G sudo -m -s /bin/bash awsuser && \
    sed -i /etc/sudoers -re 's/^%sudo.*/%sudo ALL=(ALL:ALL) NOPASSWD: ALL/g' && \
    sed -i /etc/sudoers -re 's/^root.*/root ALL=(ALL:ALL) NOPASSWD: ALL/g' && \
    sed -i /etc/sudoers -re 's/^#includedir.*/## **Removed the include directive** ##"/g' && \
    echo "awsuser ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers && \
    echo "Customized the sudoers file for passwordless access to the awsuser user!" && \
    echo "awsuser user:";  su - awsuser -c id


RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" \
    && unzip awscliv2.zip \
    && ./aws/install

RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/` curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/usr/local/bin/kubectl \
    && chmod +x /usr/local/bin/kubectl

USER awsuser
WORKDIR /home/awsuser

ENV AZURE_TENANT_ID=0ab4cbbf-4bc7-4826-b52c-a14fed5286b9
ENV APP_ID_URI=https://signin.aws.amazon.com/saml
ENV AZURE_DEFAULT_USERNAME=viacheslav.kim@rci.rogers.ca
ENV AZURE_DEFAULT_PASSWORD=gZ05QRsn6v5mfP76
ENV AZURE_DEFAULT_DURATION_HOURS=12

CMD ["/bin/bash"]